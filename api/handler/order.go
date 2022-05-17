package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type Order interface {
	ListOrder(limit, skip int) ([]byte, int)
	DetailOrder(id int64) ([]byte, int)
	SubTotalOrder(orderRequest []model.OrderedProduct) ([]byte, int)
	AddOrder(product model.AddOrderRequest) ([]byte, int)
	DownloadOrder(id int64) ([]byte, int)
	CheckOrderDownload(id int64) ([]byte, int)
}

func (s service) ListOrder(limit, skip int) ([]byte, int) {

	orders, err := s.db.GetOrder(s.ctx, limit, skip)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, orders)
}

func (s service) DetailOrder(id int64) ([]byte, int) {
	order, err := s.db.GetOrderByID(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, order)
}

func (s service) SubTotalOrder(orderRequest []model.OrderedProduct) ([]byte, int) {
	for _, item := range orderRequest {
		err := s.validation.Struct(item)
		if err != nil {
			return utils.ResponseWrapper(http.StatusBadRequest, nil)
		}
	}
	var productIds []int64
	mapProductQty := make(map[int64]int)
	for _, productItem := range orderRequest {
		productIds = append(productIds, productItem.ProductId)
		mapProductQty[productItem.ProductId] = productItem.Qty
		orderRequest = orderRequest[1:]
	}
	sort.Slice(orderRequest, func(i, j int) bool {
		return orderRequest[i].ProductId < orderRequest[j].ProductId
	})

	products, err := s.db.GetProductsByIds(s.ctx, productIds)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	orderedProductDetails, totalPrice := s.generateOrderedProduct(products, mapProductQty)
	subTotalOrder := model.SubTotalOrder{
		Subtotal:       totalPrice,
		OrderedProduct: orderedProductDetails,
	}
	return utils.ResponseWrapper(http.StatusOK, subTotalOrder)
}

func (s service) AddOrder(orderRequest model.AddOrderRequest) ([]byte, int) {

	var productIds []int64
	mapProductQty := make(map[int64]int)
	for _, productItem := range orderRequest.OrderedProduct {
		productIds = append(productIds, productItem.ProductId)
		mapProductQty[productItem.ProductId] = productItem.Qty
		orderRequest.OrderedProduct = orderRequest.OrderedProduct[1:]
	}
	sort.Slice(orderRequest.OrderedProduct, func(i, j int) bool {
		return orderRequest.OrderedProduct[i].ProductId < orderRequest.OrderedProduct[j].ProductId
	})

	products, err := s.db.GetProductsByIds(s.ctx, productIds)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	orderedProductDetails, totalPrice := s.generateOrderedProduct(products, mapProductQty)
	now := time.Now()
	order := model.Order{
		PaymentID:      orderRequest.PaymentID,
		OrderedProduct: orderedProductDetails,
		TotalPaid:      orderRequest.TotalPaid,
		TotalPrice:     totalPrice,
		TotalReturn:    orderRequest.TotalPaid - totalPrice,
		CreatedAt:      &now,
	}

	order, err = s.db.CreateOrder(s.ctx, order)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	paymentType, err := s.getPaymentType(s.ctx, orderRequest.PaymentID)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	order.PaymentType = paymentType
	return utils.ResponseWrapper(http.StatusOK, order)
}

func (s service) DownloadOrder(id int64) ([]byte, int) {
	receiptPath, err := s.db.DownloadReceipt(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, receiptPath)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, receiptPath)
	}
	return utils.ResponseWrapper(http.StatusOK, receiptPath)
}

func (s service) CheckOrderDownload(id int64) ([]byte, int) {
	isDownloaded, err := s.db.GetDownloadStatus(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	isDownloadedJson := map[string]interface{}{"isDownload": isDownloaded}
	return utils.ResponseWrapper(http.StatusOK, isDownloadedJson)
}

func (s service) calculatePrice(discount model.Discount, price, qty int) int {
	var finalPrice int
	if discount.Type == "PERCENT" {
		normalPrice := price * qty
		discountPrice := normalPrice * discount.Result / 100
		finalPrice = normalPrice - discountPrice
	} else {
		if qty >= discount.Qty {
			finalPrice = qty*price - discount.Qty*discount.Result
		}
	}

	return finalPrice
}

func (s service) generateOrderedProduct(
	products []model.Product, mapProductQty map[int64]int) ([]model.OrderedProductDetail, int) {
	var orderedProductDetails []model.OrderedProductDetail
	var totalPrice int
	for _, productItem := range products {
		productQty := mapProductQty[productItem.ProductId]
		if productQty == 0 {
			continue
		}
		var finalPrice int
		normalPrice := productItem.Price * productQty
		var discountID *int
		if productItem.Discount != nil {
			discountID = &productItem.Discount.DiscountID
			finalPrice = s.calculatePrice(*productItem.Discount,
				productItem.Price,
				productQty)
		} else {
			finalPrice = normalPrice
		}
		totalPrice += finalPrice
		orderedProductDetail := model.OrderedProductDetail{
			Product:          productItem,
			TotalFinalPrice:  finalPrice,
			TotalNormalPrice: normalPrice,
			Qty:              productQty,
			DiscountID:       discountID,
		}
		orderedProductDetails = append(orderedProductDetails, orderedProductDetail)
	}

	return orderedProductDetails, totalPrice
}

func (s service) getPaymentType(ctx context.Context, id int64) (model.Payment, error) {
	var payment model.Payment
	paymentData, err := s.db.GetPaymentByID(ctx, id)
	if err != nil {
		return payment, err
	}
	return paymentData, nil
}

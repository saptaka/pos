package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type Order interface {
	ListOrder(limit, skip int) ([]byte, int)
	DetailOrder(id int64, receiptId string) ([]byte, int)
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
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, orders)
}

func (s service) DetailOrder(id int64, receiptId string) ([]byte, int) {

	orderChan := make(chan model.Order)
	errorOrderChan := make(chan error)
	go func(orderChan chan model.Order,
		errorOrderChan chan error) {
		var order model.Order
		var err error
		if id != 0 {
			order, err = s.db.GetOrderByID(s.ctx, id)
		} else {
			order, err = s.db.GetOrderByReceiptID(s.ctx, receiptId)
		}
		if err != nil {
			orderChan <- model.Order{}
			errorOrderChan <- err
			return
		}
		orderChan <- order
		errorOrderChan <- nil
	}(orderChan, errorOrderChan)

	orderedProductChan := make(chan []model.OrderedProductDetail)
	errorOrderedProductChan := make(chan error)

	go func(orderedProductChan chan []model.OrderedProductDetail,
		errorOrderedProductChan chan error) {
		orderedProducts, err := s.db.GetOrderedProductByOrderId(s.ctx, id)
		if err != nil {
			orderedProductChan <- make([]model.OrderedProductDetail, 0)
			errorOrderedProductChan <- err
			return
		}
		orderedProductChan <- orderedProducts
		errorOrderedProductChan <- nil
	}(orderedProductChan, errorOrderedProductChan)

	order := <-orderChan
	errOrder := <-errorOrderChan
	orderedProducts := <-orderedProductChan
	errOrderedProduct := <-errorOrderedProductChan

	if errOrder == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}

	if errOrder != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}

	if errOrderedProduct == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}

	if errOrderedProduct != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}

	orders := model.Orders{
		Order:          order,
		OrderedProduct: orderedProducts,
	}

	return utils.ResponseWrapper(http.StatusOK, orders)
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
	}

	products, err := s.db.GetProductsByIds(s.ctx, productIds)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	orderedProducts, totalPrice := s.generateOrderedProduct(products, mapProductQty)
	now := time.Now()
	order := model.Order{
		PaymentID:   &orderRequest.PaymentID,
		TotalPaid:   orderRequest.TotalPaid,
		TotalPrice:  totalPrice,
		TotalReturn: orderRequest.TotalPaid - totalPrice,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		ReceiptID:   s.generateOrderID(),
	}

	order, err = s.db.CreateOrder(s.ctx, order)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}

	err = s.db.CreateOrderedProduct(context.Background(), order.OrderId, orderedProducts)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}

	orders := model.Orders{
		Order:          order,
		OrderedProduct: orderedProducts,
	}

	return utils.ResponseWrapper(http.StatusOK, orders)
}

func (s service) DownloadOrder(id int64) ([]byte, int) {
	receiptPath, err := s.db.DownloadReceipt(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, receiptPath)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, receiptPath)
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

		if productItem.Discount != nil {
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
		}
		orderedProductDetails = append(orderedProductDetails, orderedProductDetail)
	}

	return orderedProductDetails, totalPrice
}

func (s service) generateOrderID() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 1)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 999
	middleDigit := fmt.Sprint(rand.Intn(max-min+1) + min)
	return "S" + middleDigit + string(b)

}

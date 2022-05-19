package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

	orderedProductDetails, totalPrice, err := s.generateOrderedProduct(orderRequest)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	subTotalOrder := model.SubTotalOrder{
		Subtotal:       totalPrice,
		OrderedProduct: orderedProductDetails,
	}
	return utils.ResponseWrapper(http.StatusOK, subTotalOrder)
}

func (s service) AddOrder(orderRequest model.AddOrderRequest) ([]byte, int) {

	var totalPrice int
	var orderedProductDetails []model.OrderedProductDetail
	orderedProductDetails, totalPrice, err := s.generateOrderedProduct(orderRequest.OrderedProduct)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
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

	err = s.db.CreateOrderedProduct(context.Background(), order.OrderId, orderedProductDetails)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}

	orders := model.Orders{
		Order:          order,
		OrderedProduct: orderedProductDetails,
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
	orderRequest []model.OrderedProduct) ([]model.OrderedProductDetail, int, error) {
	var totalPrice int
	var orderedProductDetails []model.OrderedProductDetail
	for _, productItem := range orderRequest {

		var product model.Product
		var err error
		if productCache[productItem.ProductId] != nil {
			product = *productCache[productItem.ProductId]
		} else {
			product, err = s.db.GetProductByID(s.ctx, productItem.ProductId)
			if err != nil {
				log.Println(err)
				return orderedProductDetails, totalPrice, err
			}
		}

		if product.Stock < productItem.Qty {
			continue
		}

		go func() {
			product.Stock = product.Stock - productItem.Qty
			err := s.db.UpdateProduct(s.ctx, product)
			if err != nil {
				log.Println(err)
			}
		}()

		var finalPrice int
		normalPrice := product.Price * productItem.Qty
		if product.DiscountId != nil {
			finalPrice = s.calculatePrice(*product.Discount,
				product.Price,
				productItem.Qty)
		} else {
			finalPrice = normalPrice
		}

		totalPrice += finalPrice
		orderedProductDetail := model.OrderedProductDetail{
			Product:          product,
			TotalFinalPrice:  finalPrice,
			TotalNormalPrice: normalPrice,
			Qty:              productItem.Qty,
		}
		orderedProductDetails = append(orderedProductDetails, orderedProductDetail)

	}

	return orderedProductDetails, totalPrice, nil
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

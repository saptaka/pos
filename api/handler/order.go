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
	"github.com/valyala/fasthttp"
)

type Order interface {
	ListOrder(limit, skip int) (map[string]interface{}, int)
	DetailOrder(id int64, receiptId string) (map[string]interface{}, int)
	SubTotalOrder(orderRequest model.OrderedProducts) (map[string]interface{}, int)
	AddOrder(product model.AddOrderRequest) (map[string]interface{}, int)
	DownloadOrder(id int64) (map[string]interface{}, int)
	CheckOrderDownload(id int64) (map[string]interface{}, int)
}

func (s service) ListOrder(limit, skip int) (map[string]interface{}, int) {

	orders, err := s.db.GetOrder(s.ctx, limit, skip)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	listOrders := model.ListOrders{
		Order: orders,
		Meta: model.Meta{
			Limit: limit,
			Skip:  skip,
			Total: len(orders),
		},
	}

	return utils.ResponseWrapper(http.StatusOK, listOrders, nil)
}

func (s service) DetailOrder(id int64, receiptId string) (map[string]interface{}, int) {

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
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}

	if errOrder != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}

	if errOrderedProduct == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}

	if errOrderedProduct != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}

	orderDetails := model.OrderDetails{
		Order:          order,
		OrderedProduct: orderedProducts,
	}

	return utils.ResponseWrapper(http.StatusOK, orderDetails, nil)
}

func (s service) SubTotalOrder(orderRequest model.OrderedProducts) (map[string]interface{}, int) {
	validation := orderValidation(model.SUBORDER)
	err := validation.Struct(orderRequest)
	if err != nil {
		return utils.ErrorWrapper(err, fasthttp.StatusBadRequest, model.SUBORDER)
	}
	orderedProductDetails, totalPrice, err := s.generateSubOrderedProduct(orderRequest.Products)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	subTotalOrder := model.SubTotalOrder{
		Subtotal:       totalPrice,
		OrderedProduct: orderedProductDetails,
	}
	return utils.ResponseWrapper(http.StatusOK, subTotalOrder, nil)
}

func (s service) AddOrder(orderRequest model.AddOrderRequest) (map[string]interface{}, int) {
	validation := orderValidation(model.CREATE)
	err := validation.Struct(orderRequest)
	if err != nil {
		return utils.ErrorWrapper(err, fasthttp.StatusBadRequest, model.CREATE)
	}
	var totalPrice int
	subOrderedProductDetails, totalPrice, err := s.generateSubOrderedProduct(orderRequest.OrderedProduct)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}

	now, _ := time.Parse(model.RFC3339MilliZ, time.Now().UTC().Format(model.RFC3339MilliZ))
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
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	var orderedProductDetails []model.OrderedProductDetail
	for _, subOderedProductDetail := range subOrderedProductDetails {
		orderedProductDetail := model.OrderedProductDetail{
			ProductId:        subOderedProductDetail.ProductId,
			Name:             subOderedProductDetail.Name,
			Price:            subOderedProductDetail.Price,
			Qty:              subOderedProductDetail.Qty,
			Discount:         subOderedProductDetail.Discount,
			TotalFinalPrice:  subOderedProductDetail.TotalFinalPrice,
			TotalNormalPrice: subOderedProductDetail.TotalNormalPrice,
		}
		orderedProductDetails = append(orderedProductDetails, orderedProductDetail)
	}

	orders := model.OrderDetails{
		Order:          order,
		OrderedProduct: orderedProductDetails,
	}

	go func() {
		err = s.db.CreateOrderedProduct(context.Background(), order.OrderId, orderedProductDetails)
		if err != nil {
			log.Println(err)
		}
	}()

	return utils.ResponseWrapper(http.StatusOK, orders, nil)
}

func (s service) DownloadOrder(id int64) (map[string]interface{}, int) {
	receiptPath, err := s.db.DownloadReceipt(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, receiptPath, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, receiptPath, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, receiptPath, nil)
}

func (s service) CheckOrderDownload(id int64) (map[string]interface{}, int) {
	isDownloaded, err := s.db.GetDownloadStatus(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	isDownloadedJson := map[string]interface{}{"isDownload": isDownloaded}
	return utils.ResponseWrapper(http.StatusOK, isDownloadedJson, nil)
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

func (s service) generateSubOrderedProduct(
	orderRequest []model.OrderedProduct) ([]model.SubOrderedProductDetail, int, error) {
	var totalPrice int
	var orderedProductDetails []model.SubOrderedProductDetail
	mapOrderedProduct := make(map[int64]int)
	for index, productItem := range orderRequest {

		var product model.Product
		var err error
		productMapValue, ok := productCache.Get(productItem.ProductId)
		if ok {
			product = productMapValue
		} else {
			product, err = s.db.GetProductByID(s.ctx, productItem.ProductId)
			if err == sql.ErrNoRows {
				continue
			}
			if err != nil {
				log.Println(err)
				return orderedProductDetails, totalPrice, err
			}
		}

		if product.Stock < productItem.Qty {
			continue
		}
		product.Stock = product.Stock - productItem.Qty

		err = s.db.UpdateProduct(s.ctx, product)
		if err != nil {
			log.Printf("error update product in order process %d : %s",
				product.ProductId, err)
		}

		productCache.Set(product.ProductId, product)

		var finalPrice int
		normalPrice := product.Price * productItem.Qty
		if product.DiscountId != nil {
			finalPrice = s.calculatePrice(*product.Discount,
				product.Price,
				productItem.Qty)
		} else {
			finalPrice = normalPrice
		}

		if orderIndex, ok := mapOrderedProduct[product.ProductId]; ok {
			orderedProductDetails[orderIndex].Qty += productItem.Qty
			orderedProductDetails[orderIndex].TotalFinalPrice += finalPrice
			orderedProductDetails[orderIndex].TotalNormalPrice += normalPrice
			orderedProductDetails[orderIndex].Stock = product.Stock
			totalPrice += finalPrice
			continue
		}

		var discount *model.Discount
		if product.Discount != nil {
			discount = product.Discount
		}

		totalPrice += finalPrice
		orderedProductDetail := model.SubOrderedProductDetail{
			Product: model.Product{
				ProductId: product.ProductId,
				Name:      product.Name,
				Price:     product.Price,
				Discount:  discount,
				Stock:     product.Stock,
				Image:     product.Image,
			},
			Qty:              productItem.Qty,
			TotalFinalPrice:  finalPrice,
			TotalNormalPrice: normalPrice,
		}
		orderedProductDetails = append(orderedProductDetails, orderedProductDetail)
		mapOrderedProduct[product.ProductId] = index
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

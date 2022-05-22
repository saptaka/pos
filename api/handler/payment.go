package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type Payment interface {
	ListPayment(limit, skip int) (map[string]interface{}, int)
	DetailPayment(id int64) (map[string]interface{}, int)
	CreatePayment(payment model.Payment) (map[string]interface{}, int)
	UpdatePayment(payment model.Payment) (map[string]interface{}, int)
	DeletePayment(id int) (map[string]interface{}, int)
}

func (s service) ListPayment(limit, skip int) (map[string]interface{}, int) {
	Payments, err := s.db.GetPayments(context.Background(), limit, skip)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	listPayment := model.ListPayment{
		Payments: Payments,
		Meta: model.Meta{
			Total: len(Payments),
			Limit: limit,
			Skip:  skip,
		},
	}

	return utils.ResponseWrapper(http.StatusOK, listPayment, nil)
}

func (s service) DetailPayment(id int64) (map[string]interface{}, int) {
	payment, err := s.db.GetPaymentByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, payment, nil)
}

func (s service) CreatePayment(payment model.Payment) (map[string]interface{}, int) {
	validation := paymentValidation(model.CREATE)
	err := validation.Struct(payment)
	if err != nil {
		return utils.ErrorWrapper(err, fasthttp.StatusBadRequest, model.CREATE)
	}

	paymentData, err := s.db.CreatePayment(s.ctx, payment)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, paymentData, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, paymentData, nil)
}

func (s service) UpdatePayment(payment model.Payment) (map[string]interface{}, int) {
	validation := paymentValidation(model.UPDATE)
	err := validation.Struct(payment)
	if err != nil {
		return utils.ErrorWrapper(err, fasthttp.StatusBadRequest, model.UPDATE)
	}

	err = s.db.UpdatePayment(s.ctx, payment)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, "Payment Not Found", nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

func (s service) DeletePayment(id int) (map[string]interface{}, int) {
	err := s.db.DeletePayment(s.ctx, id)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

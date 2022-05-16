package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type Payment interface {
	ListPayment(limit, skip int) ([]byte, int)
	DetailPayment(id int64) ([]byte, int)
	CreatePayment(payment model.Payment) ([]byte, int)
	UpdatePayment(payment model.Payment) ([]byte, int)
	DeletePayment(id int) ([]byte, int)
}

func (s service) ListPayment(limit, skip int) ([]byte, int) {
	Payments, err := s.db.GetPayments(context.Background(), limit, skip)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	listPayment := model.ListPayment{
		Payments: Payments,
		Meta: model.Meta{
			Total: len(Payments),
			Limit: limit,
			Skip:  skip,
		},
	}

	return utils.ResponseWrapper(http.StatusOK, listPayment)
}

func (s service) DetailPayment(id int64) ([]byte, int) {
	payment, err := s.db.GetPaymentByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, payment)
}

func (s service) CreatePayment(payment model.Payment) ([]byte, int) {
	err := s.validation.Struct(payment)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if !model.PaymentType[payment.Type] {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}

	paymentData, err := s.db.CreatePayment(s.ctx, payment)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, paymentData)
	}
	return utils.ResponseWrapper(http.StatusOK, paymentData)
}

func (s service) UpdatePayment(payment model.Payment) ([]byte, int) {
	err := s.db.UpdatePayment(s.ctx, payment)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

func (s service) DeletePayment(id int) ([]byte, int) {
	err := s.db.DeletePayment(s.ctx, id)
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

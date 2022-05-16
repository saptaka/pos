package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type Cashier interface {
	ListCashier(limit, skip int) ([]byte, int)
	DetailCashier(id int64) ([]byte, int)
	CreateCashier(cashier model.Cashier) ([]byte, int)
	UpdateCashier(cashier model.Cashier) ([]byte, int)
	DeleteCashier(id int64) ([]byte, int)
}

func (s service) ListCashier(limit, skip int) ([]byte, int) {
	cashiers, err := s.db.GetCashiers(context.Background(), limit, skip)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	listCashier := model.ListCashier{
		Cashiers: cashiers,
		Meta: model.Meta{
			Total: len(cashiers),
			Limit: limit,
			Skip:  skip,
		},
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, listCashier)
	}
	return utils.ResponseWrapper(http.StatusOK, listCashier)
}

func (s service) DetailCashier(id int64) ([]byte, int) {
	cashier, err := s.db.GetCashierByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, cashier)
}

func (s service) CreateCashier(cashierDetail model.Cashier) ([]byte, int) {
	err := s.validation.Struct(cashierDetail)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	_, err = strconv.Atoi(cashierDetail.Passcode)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	cashier, err := s.db.CreateCashier(s.ctx, cashierDetail.Name, cashierDetail.Passcode)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, cashier)
	}
	return utils.ResponseWrapper(http.StatusOK, cashier)
}

func (s service) UpdateCashier(cashierDetail model.Cashier) ([]byte, int) {

	err := s.db.UpdateCashier(s.ctx, cashierDetail)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

func (s service) DeleteCashier(id int64) ([]byte, int) {
	err := s.db.DeleteCashier(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

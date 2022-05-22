package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type Cashier interface {
	ListCashier(limit, skip int) (map[string]interface{}, int)
	DetailCashier(id int64) (map[string]interface{}, int)
	CreateCashier(cashier model.Cashier) (map[string]interface{}, int)
	UpdateCashier(cashier model.Cashier) (map[string]interface{}, int)
	DeleteCashier(id int64) (map[string]interface{}, int)
}

func (s service) ListCashier(limit, skip int) (map[string]interface{}, int) {
	cashiers, err := s.db.GetCashiers(context.Background(), limit, skip)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	listCashier := model.ListCashier{
		Cashiers: cashiers,
		Meta: model.Meta{
			Total: len(cashiers),
			Limit: limit,
			Skip:  skip,
		},
	}

	return utils.ResponseWrapper(http.StatusOK, listCashier, nil)
}

func (s service) DetailCashier(id int64) (map[string]interface{}, int) {
	cashier, err := s.db.GetCashierByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, cashier, nil)
}

func (s service) CreateCashier(cashierDetail model.Cashier) (map[string]interface{}, int) {
	validation := generateCashierCreateValidation(s.validation)
	err := validation.Struct(cashierDetail)
	if err != nil {
		return utils.ErrorWrapper(err, fasthttp.StatusBadRequest)
	}

	_, err = strconv.Atoi(cashierDetail.Passcode)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	cashier, err := s.db.CreateCashier(s.ctx, cashierDetail.Name, cashierDetail.Passcode)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, cashier, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, cashier, nil)
}

func (s service) UpdateCashier(cashierDetail model.Cashier) (map[string]interface{}, int) {

	err := s.db.UpdateCashier(s.ctx, cashierDetail)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

func (s service) DeleteCashier(id int64) (map[string]interface{}, int) {
	err := s.db.DeleteCashier(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

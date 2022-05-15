package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type CashierRouter interface {
	ListCashier(res http.ResponseWriter, req *http.Request)
	DetailCashier(res http.ResponseWriter, req *http.Request)
	CreateCashier(res http.ResponseWriter, req *http.Request)
	UpdateCashier(res http.ResponseWriter, req *http.Request)
	DeleteCashier(res http.ResponseWriter, req *http.Request)
	RouteCashierPath()
}

func (r *router) RouteCashierPath() {
	r.mux.HandleFunc("/cashiers", verifyToken(r.ListCashier)).Methods("GET")
	r.mux.HandleFunc("/cashiers/{cashierId}", verifyToken(r.DetailCashier)).Methods("GET")
	r.mux.HandleFunc("/cashiers", verifyToken(r.CreateCashier)).Methods("POST")
	r.mux.HandleFunc("/cashiers/{cashierId}", verifyToken(r.UpdateCashier)).Methods("PUT")
	r.mux.HandleFunc("/cashiers/{cashierId}", verifyToken(r.DeleteCashier)).Methods("DELETE")
}

func (r *router) ListCashier(res http.ResponseWriter, req *http.Request) {

	limitQuery := req.URL.Query().Get("limit")
	skipQuery := req.URL.Query().Get("skip")
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)
	response, statusCode := r.handlerService.ListCashier(limit, skip)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DetailCashier(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}
	response, statusCode := r.handlerService.DetailCashier(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) CreateCashier(res http.ResponseWriter, req *http.Request) {

	var cashierRequest model.Cashier
	err := json.NewDecoder(req.Body).Decode(&cashierRequest)
	if err != nil {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}

	response, statusCode := r.handlerService.CreateCashier(cashierRequest)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) UpdateCashier(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.Atoi(idParams)
	if id == 0 {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}

	var cashierDetail model.Cashier
	err := json.NewDecoder(req.Body).Decode(&cashierDetail)
	if err != nil {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}
	cashierDetail.ChashierId = int64(id)
	response, statusCode := r.handlerService.UpdateCashier(cashierDetail)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DeleteCashier(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}
	response, statusCode := r.handlerService.DeleteCashier(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

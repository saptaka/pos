package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type PaymentRouter interface {
	ListPayment(res http.ResponseWriter, req *http.Request)
	DetailPayment(res http.ResponseWriter, req *http.Request)
	CreatePayment(res http.ResponseWriter, req *http.Request)
	UpdatePayment(res http.ResponseWriter, req *http.Request)
	DeletePayment(res http.ResponseWriter, req *http.Request)
	RoutePaymentPath()
}

func (r *router) RoutePaymentPath() {
	r.mux.HandleFunc("/payments", middleware(r.ListPayment)).Methods("GET")
	r.mux.HandleFunc("/payments/{paymentId}", middleware(r.DetailPayment)).Methods("GET")
	r.mux.HandleFunc("/payments", r.CreatePayment).Methods("POST")
	r.mux.HandleFunc("/payments/{paymentId}", r.UpdatePayment).Methods("PUT")
	r.mux.HandleFunc("/payments/{paymentId}", r.DeletePayment).Methods("DELETE")
}

func (r *router) ListPayment(res http.ResponseWriter, req *http.Request) {

	limitQuery := req.URL.Query().Get("limit")
	skipQuery := req.URL.Query().Get("skip")
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)
	response, statusCode := r.handlerService.ListPayment(limit, skip)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DetailPayment(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["paymentId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DetailPayment(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) CreatePayment(res http.ResponseWriter, req *http.Request) {
	var payment model.Payment
	err := json.NewDecoder(req.Body).Decode(&payment)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}

	response, statusCode := r.handlerService.CreatePayment(payment)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) UpdatePayment(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["paymentId"]
	id, _ := strconv.Atoi(idParams)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}

	var payment model.Payment
	err := json.NewDecoder(req.Body).Decode(&payment)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	payment.PaymentId = int64(id)
	response, statusCode := r.handlerService.UpdatePayment(payment)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DeletePayment(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["paymentId"]
	id, _ := strconv.Atoi(idParams)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DeletePayment(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

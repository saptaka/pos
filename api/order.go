package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type OrderRouter interface {
	ListOrder(res http.ResponseWriter, req *http.Request)
	DetailOrder(res http.ResponseWriter, req *http.Request)
	SubTotalOrder(res http.ResponseWriter, req *http.Request)
	AddOrder(res http.ResponseWriter, req *http.Request)
	DownloadOrder(res http.ResponseWriter, req *http.Request)
	CheckOrderDownload(res http.ResponseWriter, req *http.Request)
	RouteOrderPath()
}

func (r *router) RouteOrderPath() {
	r.mux.HandleFunc("/orders", middleware(r.ListOrder)).Methods("GET")
	r.mux.HandleFunc("/orders/{orderId}", middleware(r.DetailOrder)).Methods("GET")
	r.mux.HandleFunc("/orders/subtotal", middleware(r.SubTotalOrder)).Methods("POST")
	r.mux.HandleFunc("/orders", middleware(r.AddOrder)).Methods("POST")
	r.mux.HandleFunc("/orders/{orderId}/download", middleware(r.DownloadOrder)).Methods("GET")
	r.mux.HandleFunc("/orders/{orderId}", middleware(r.CheckOrderDownload)).Methods("GET")
}

func (r *router) ListOrder(res http.ResponseWriter, req *http.Request) {
	limitQuery := req.URL.Query().Get("limit")
	skipQuery := req.URL.Query().Get("skip")
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)

	response, statusCode := r.handlerService.ListOrder(limit, skip)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DetailOrder(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DetailOrder(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) SubTotalOrder(res http.ResponseWriter, req *http.Request) {
	var orderedProducts []model.OrderedProduct
	err := json.NewDecoder(req.Body).Decode(&orderedProducts)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.SubTotalOrder(orderedProducts)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) AddOrder(res http.ResponseWriter, req *http.Request) {

	var addOrderRequest model.AddOrderRequest
	err := json.NewDecoder(req.Body).Decode(&addOrderRequest)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.AddOrder(addOrderRequest)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DownloadOrder(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DownloadOrder(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) CheckOrderDownload(res http.ResponseWriter, req *http.Request) {
	response, statusCode := r.handlerService.CheckOrderDownload()
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

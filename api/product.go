package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type ProductRouter interface {
	ListProduct(res http.ResponseWriter, req *http.Request)
	DetailProduct(res http.ResponseWriter, req *http.Request)
	CreateProduct(res http.ResponseWriter, req *http.Request)
	UpdateProduct(res http.ResponseWriter, req *http.Request)
	DeleteProduct(res http.ResponseWriter, req *http.Request)
	RouteProductPath()
}

func (r *router) RouteProductPath() {
	r.mux.HandleFunc("/products", middleware(r.ListProduct)).Methods("GET")
	r.mux.HandleFunc("/products/{productId}", middleware(r.DetailProduct)).Methods("GET")
	r.mux.HandleFunc("/products", r.CreateProduct).Methods("POST")
	r.mux.HandleFunc("/products/{productId}", (r.UpdateProduct)).Methods("PUT")
	r.mux.HandleFunc("/products/{productId}", r.DeleteProduct).Methods("DELETE")
}

func (r *router) ListProduct(res http.ResponseWriter, req *http.Request) {

	limitQuery := req.URL.Query().Get("limit")
	skipQuery := req.URL.Query().Get("skip")
	query := req.URL.Query().Get("q")
	categoryIdParams := req.URL.Query().Get("categoryId")
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)
	categoryId, _ := strconv.ParseInt(categoryIdParams, 10, 0)
	response, statusCode := r.handlerService.ListProduct(limit, skip, categoryId, query)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DetailProduct(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["productId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DetailProduct(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) CreateProduct(res http.ResponseWriter, req *http.Request) {

	var product model.ProductCreateRequest
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}

	response, statusCode := r.handlerService.CreateProduct(product)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) UpdateProduct(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["productId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	var product model.Product
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	product.ProductId = id
	response, statusCode := r.handlerService.UpdateProduct(product)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DeleteProduct(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["productId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DeleteProduct(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

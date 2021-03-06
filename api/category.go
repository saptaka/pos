package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type CategoryRouter interface {
	ListCategory(res http.ResponseWriter, req *http.Request)
	DetailCategory(res http.ResponseWriter, req *http.Request)
	CreateCategory(res http.ResponseWriter, req *http.Request)
	UpdateCategory(res http.ResponseWriter, req *http.Request)
	DeleteCategory(res http.ResponseWriter, req *http.Request)
	RouteCategoryPath()
}

func (r *router) RouteCategoryPath() {
	r.mux.HandleFunc("/categories", middleware(r.ListCategory)).Methods("GET")
	r.mux.HandleFunc("/categories/{categoryId}", middleware(r.DetailCategory)).Methods("GET")
	r.mux.HandleFunc("/categories", r.CreateCategory).Methods("POST")
	r.mux.HandleFunc("/categories/{categoryId}", r.UpdateCategory).Methods("PUT")
	r.mux.HandleFunc("/categories/{categoryId}", r.DeleteCategory).Methods("DELETE")
}

func (r *router) ListCategory(res http.ResponseWriter, req *http.Request) {

	limitQuery := req.URL.Query().Get("limit")
	skipQuery := req.URL.Query().Get("skip")
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)
	response, statusCode := r.handlerService.ListCategory(limit, skip)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DetailCategory(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["categoryId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DetailCategory(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) CreateCategory(res http.ResponseWriter, req *http.Request) {

	var category model.Category
	err := json.NewDecoder(req.Body).Decode(&category)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}

	response, statusCode := r.handlerService.CreateCategory(category)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) UpdateCategory(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["categoryId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}

	var category model.Category
	err := json.NewDecoder(req.Body).Decode(&category)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	category.CategoryId = id
	response, statusCode := r.handlerService.UpdateCategory(category)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) DeleteCategory(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["categoryId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil)
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	response, statusCode := r.handlerService.DeleteCategory(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write(response)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

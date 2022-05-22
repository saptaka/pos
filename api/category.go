package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type CategoryRouter interface {
	ListCategory(req *fasthttp.RequestCtx)
	DetailCategory(req *fasthttp.RequestCtx)
	CreateCategory(req *fasthttp.RequestCtx)
	UpdateCategory(req *fasthttp.RequestCtx)
	DeleteCategory(req *fasthttp.RequestCtx)
	RouteCategoryPath()
}

func (r *apiRouter) RouteCategoryPath() {
	r.mux.GET("/categories", middleware(r.ListCategory))
	r.mux.GET("/categories/{categoryId}", middleware(r.DetailCategory))
	r.mux.POST("/categories", r.CreateCategory)
	r.mux.PUT("/categories/{categoryId}", r.UpdateCategory)
	r.mux.DELETE("/categories/{categoryId}", r.DeleteCategory)
}

func (r *apiRouter) ListCategory(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	limitQuery := string(req.URI().QueryArgs().Peek("limit"))
	skipQuery := string(req.URI().QueryArgs().Peek("skip"))
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)
	response, statusCode := r.handlerService.ListCategory(limit, skip)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DetailCategory(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("categoryId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DetailCategory(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) CreateCategory(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	var category model.Category
	err := json.Unmarshal(req.Request.Body(), &category)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	response, statusCode := r.handlerService.CreateCategory(category)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) UpdateCategory(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("categoryId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	var category model.Category
	err := json.Unmarshal(req.Request.Body(), &category)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	category.CategoryId = id
	response, statusCode := r.handlerService.UpdateCategory(category)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DeleteCategory(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("categoryId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DeleteCategory(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

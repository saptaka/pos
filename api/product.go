package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type ProductRouter interface {
	ListProduct(req *fasthttp.RequestCtx)
	DetailProduct(req *fasthttp.RequestCtx)
	CreateProduct(req *fasthttp.RequestCtx)
	UpdateProduct(req *fasthttp.RequestCtx)
	DeleteProduct(req *fasthttp.RequestCtx)
	RouteProductPath()
}

func (r *apiRouter) RouteProductPath() {
	r.mux.GET("/products", middleware(r.ListProduct))
	r.mux.GET("/products/{productId}", middleware(r.DetailProduct))
	r.mux.POST("/products", r.CreateProduct)
	r.mux.PUT("/products/{productId}", (r.UpdateProduct))
	r.mux.DELETE("/products/{productId}", r.DeleteProduct)
}

func (r *apiRouter) ListProduct(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	limitQuery := string(req.URI().QueryArgs().Peek("limit"))
	skipQuery := string(req.URI().QueryArgs().Peek("skip"))
	query := string(req.URI().QueryArgs().Peek("q"))
	categoryIdParams := string(req.URI().QueryArgs().Peek("categoryId"))
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)
	categoryId, err := strconv.ParseInt(categoryIdParams, 10, 0)
	var product model.Product
	if err == nil {
		product = model.Product{
			CategoryId: &categoryId,
		}
	}

	if query != "" {
		product.Name = query
	}
	response, statusCode := r.handlerService.ListProduct(limit, skip, product)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DetailProduct(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	idParams := req.UserValue("productId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DetailProduct(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) CreateProduct(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	var product model.ProductCreateRequest
	err := json.Unmarshal(req.Request.Body(), &product)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	response, statusCode := r.handlerService.CreateProduct(product)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) UpdateProduct(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	idParams := req.UserValue("productId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	var product model.Product
	err := json.Unmarshal(req.Request.Body(), &product)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	product.ProductId = id
	response, statusCode := r.handlerService.UpdateProduct(product)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DeleteProduct(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	idParams := req.UserValue("productId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DeleteProduct(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

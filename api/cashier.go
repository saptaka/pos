package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type CashierRouter interface {
	ListCashier(req *fasthttp.RequestCtx)
	DetailCashier(req *fasthttp.RequestCtx)
	CreateCashier(req *fasthttp.RequestCtx)
	UpdateCashier(req *fasthttp.RequestCtx)
	DeleteCashier(req *fasthttp.RequestCtx)
	RouteCashierPath()
}

func (r *apiRouter) RouteCashierPath() {
	r.mux.GET("/cashiers", r.ListCashier)
	r.mux.GET("/cashiers/{cashierId}", r.DetailCashier)
	r.mux.POST("/cashiers", r.CreateCashier)
	r.mux.PUT("/cashiers/{cashierId}", r.UpdateCashier)
	r.mux.DELETE("/cashiers/{cashierId}", r.DeleteCashier)
}

func (r *apiRouter) ListCashier(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	limitQuery := req.URI().QueryArgs().Peek("limit")
	skipQuery := req.URI().QueryArgs().Peek("skip")
	limit, _ := strconv.Atoi(string(limitQuery))
	skip, _ := strconv.Atoi(string(skipQuery))
	response, statusCode := r.handlerService.ListCashier(limit, skip)
	if statusCode != http.StatusOK {
		response, statusCode := utils.ResponseWrapper(fasthttp.StatusBadRequest, "")
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DetailCashier(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("cashierId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DetailCashier(id)
	if statusCode != http.StatusOK {

		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) CreateCashier(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	var cashierRequest model.Cashier
	err := json.Unmarshal(req.Request.Body(), &cashierRequest)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	response, statusCode := r.handlerService.CreateCashier(cashierRequest)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) UpdateCashier(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("cashierId").(string)
	id, _ := strconv.Atoi(idParams)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil)
		req.Response.Header.SetCanonical(model.ContentTypeJSON())
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	var cashierDetail model.Cashier
	err := json.Unmarshal(req.Request.Body(), &cashierDetail)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.Header.SetCanonical(model.ContentTypeJSON())
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	cashierDetail.CashierId = int64(id)
	response, statusCode := r.handlerService.UpdateCashier(cashierDetail)
	if statusCode != http.StatusOK {
		req.Response.Header.SetCanonical(model.ContentTypeJSON())
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DeleteCashier(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("cashierId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DeleteCashier(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

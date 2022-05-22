package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type PaymentRouter interface {
	ListPayment(req *fasthttp.RequestCtx)
	DetailPayment(req *fasthttp.RequestCtx)
	CreatePayment(req *fasthttp.RequestCtx)
	UpdatePayment(req *fasthttp.RequestCtx)
	DeletePayment(req *fasthttp.RequestCtx)
	RoutePaymentPath()
}

func (r *apiRouter) RoutePaymentPath() {
	r.mux.GET("/payments", middleware(r.ListPayment))
	r.mux.GET("/payments/{paymentId}", middleware(r.DetailPayment))
	r.mux.POST("/payments", r.CreatePayment)
	r.mux.PUT("/payments/{paymentId}", r.UpdatePayment)
	r.mux.DELETE("/payments/{paymentId}", r.DeletePayment)
}

func (r *apiRouter) ListPayment(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	limitQuery := string(req.URI().QueryArgs().Peek("limit"))
	skipQuery := string(req.URI().QueryArgs().Peek("skip"))
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)
	response, statusCode := r.handlerService.ListPayment(limit, skip)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DetailPayment(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	idParams := req.UserValue("paymentId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DetailPayment(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) CreatePayment(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	var payment model.Payment
	err := json.Unmarshal(req.Request.Body(), &payment)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	response, statusCode := r.handlerService.CreatePayment(payment)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) UpdatePayment(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	idParams := req.UserValue("paymentId").(string)
	id, _ := strconv.Atoi(idParams)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusNotFound, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	var payment model.Payment
	err := json.Unmarshal(req.Request.Body(), &payment)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	payment.PaymentId = int64(id)
	response, statusCode := r.handlerService.UpdatePayment(payment)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DeletePayment(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	idParams := req.UserValue("paymentId").(string)
	id, _ := strconv.Atoi(idParams)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DeletePayment(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

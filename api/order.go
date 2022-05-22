package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type OrderRouter interface {
	ListOrder(req *fasthttp.RequestCtx)
	DetailOrder(req *fasthttp.RequestCtx)
	SubTotalOrder(req *fasthttp.RequestCtx)
	AddOrder(req *fasthttp.RequestCtx)
	DownloadOrder(req *fasthttp.RequestCtx)
	CheckOrderDownload(req *fasthttp.RequestCtx)
	RouteOrderPath()
}

func (r *apiRouter) RouteOrderPath() {
	r.mux.GET("/orders", r.ListOrder)
	r.mux.GET("/orders/{orderId}", middleware(r.DetailOrder))
	r.mux.POST("/orders/subtotal", middleware(r.SubTotalOrder))
	r.mux.POST("/orders", middleware(r.AddOrder))
	r.mux.GET("/orders/{orderId}/download", middleware(r.DownloadOrder))
	r.mux.GET("/orders/{orderId}/check-download", middleware(r.CheckOrderDownload))
}

func (r *apiRouter) ListOrder(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	limitQuery := string(req.URI().QueryArgs().Peek("limit"))
	skipQuery := string(req.URI().QueryArgs().Peek("skip"))
	limit, _ := strconv.Atoi(limitQuery)
	skip, _ := strconv.Atoi(skipQuery)

	response, statusCode := r.handlerService.ListOrder(limit, skip)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DetailOrder(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("orderId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	re, err := regexp.Compile(`S[0-9]{3}[A-Z]{1}`)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	isReceiptId := re.MatchString(idParams)
	if id == 0 && !isReceiptId {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	var receiptId string
	if isReceiptId {
		receiptId = idParams
	}

	response, statusCode := r.handlerService.DetailOrder(id, receiptId)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) SubTotalOrder(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	var orderedProducts []model.OrderedProduct
	err := json.Unmarshal(req.Request.Body(), &orderedProducts)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	if len(orderedProducts) == 0 {
		response, _ := utils.ResponseWrapper(http.StatusBadRequest, "\"value\" must be an array",
			[]model.ErrorData{{
				Message: "\"value\" must be an array",
				Path:    []string{},
				Type:    "array.base",
				Context: model.CreateErrorContext{
					Label: "value",
					Key:   make(map[string]interface{}),
				},
			}})
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.SubTotalOrder(orderedProducts)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) AddOrder(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	var addOrderRequest model.AddOrderRequest
	err := json.Unmarshal(req.Request.Body(), &addOrderRequest)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.AddOrder(addOrderRequest)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) DownloadOrder(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("orderId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.DownloadOrder(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) CheckOrderDownload(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("orderId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.CheckOrderDownload(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

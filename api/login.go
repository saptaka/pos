package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type LoginRouter interface {
	GetPasscode(req *fasthttp.RequestCtx)
	VerifyLogin(req *fasthttp.RequestCtx)
	VerifyLogout(req *fasthttp.RequestCtx)
	RouteLoginPath()
}

func (r *apiRouter) RouteLoginPath() {
	r.mux.GET("/cashiers/{cashierId}/passcode", r.GetPasscode)
	r.mux.POST("/cashiers/{cashierId}/login", r.VerifyLogin)
	r.mux.POST("/cashiers/{cashierId}/logout", r.VerifyLogout)
}

func (r *apiRouter) GetPasscode(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := string(req.URI().QueryArgs().Peek("cashierId"))
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}

	response, statusCode := r.handlerService.GetPasscode(id)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) VerifyLogin(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("cashierId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	var cashier model.Cashier
	err := json.Unmarshal(req.Request.Body(), &cashier)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.VerifyLogin(id, cashier.Passcode, Token)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) VerifyLogout(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	idParams := req.UserValue("cashierId").(string)
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	var cashier model.Cashier
	err := json.Unmarshal(req.Request.Body(), &cashier)
	if err != nil {
		response, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	response, statusCode := r.handlerService.VerifyLogout(id, cashier.Passcode)
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

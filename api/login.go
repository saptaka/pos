package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type LoginRouter interface {
	GetPasscode(res http.ResponseWriter, req *http.Request)
	VerifyLogin(res http.ResponseWriter, req *http.Request)
	VerifyLogout(res http.ResponseWriter, req *http.Request)
	RouteLoginPath()
}

func (r *router) RouteLoginPath() {
	r.mux.HandleFunc("/cashiers/{cashierId}/passcode", r.GetPasscode).Methods("GET")
	r.mux.HandleFunc("/cashiers/{cashierId}/login", r.VerifyLogin).Methods("POST")
	r.mux.HandleFunc("/cashiers/{cashierId}/logout", r.VerifyLogout).Methods("POST")
}

func (r *router) GetPasscode(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}

	response, statusCode := r.handlerService.GetPasscode(id)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) VerifyLogin(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}
	var cashier model.Cashier
	err := json.NewDecoder(req.Body).Decode(&cashier)
	if err != nil {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}
	response, statusCode := r.handlerService.VerifyLogin(id, cashier.Passcode, Token)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) VerifyLogout(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	idParams := params["cashierId"]
	id, _ := strconv.ParseInt(idParams, 10, 0)
	if id == 0 {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}
	var cashier model.Cashier
	err := json.NewDecoder(req.Body).Decode(&cashier)
	if err != nil {
		_, statusCode := utils.ResponseWrapper(http.StatusBadRequest, nil)
		res.WriteHeader(statusCode)
		return
	}
	response, statusCode := r.handlerService.VerifyLogout(id, cashier.Passcode)
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

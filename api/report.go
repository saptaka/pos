package api

import (
	"net/http"
)

type ReportRouter interface {
	Revenue(res http.ResponseWriter, req *http.Request)
	Solds(res http.ResponseWriter, req *http.Request)
	RouteReportPath()
}

func (r *router) RouteReportPath() {
	r.mux.HandleFunc("/revenues", verifyToken(r.Revenue)).Methods("GET")
	r.mux.HandleFunc("/solds", verifyToken(r.Solds)).Methods("GET")
}

func (r *router) Revenue(res http.ResponseWriter, req *http.Request) {

	response, statusCode := r.handlerService.Revenue()
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

func (r *router) Solds(res http.ResponseWriter, req *http.Request) {
	response, statusCode := r.handlerService.Solds()
	if statusCode != http.StatusOK {
		res.WriteHeader(statusCode)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(response)
}

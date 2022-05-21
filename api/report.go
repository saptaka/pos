package api

import (
	"encoding/json"
	"net/http"

	"github.com/saptaka/pos/model"
	"github.com/valyala/fasthttp"
)

type ReportRouter interface {
	Revenue(req *fasthttp.RequestCtx)
	Solds(req *fasthttp.RequestCtx)
	RouteReportPath()
}

func (r *apiRouter) RouteReportPath() {
	r.mux.GET("/revenues", middleware(r.Revenue))
	r.mux.GET("/solds", middleware(r.Solds))
}

func (r *apiRouter) Revenue(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())

	response, statusCode := r.handlerService.Revenue()
	if statusCode != http.StatusOK {
		req.Response.SetStatusCode(statusCode)
		return
	}
	json.NewEncoder(req).Encode(response)
}

func (r *apiRouter) Solds(req *fasthttp.RequestCtx) {
	req.Response.Header.SetCanonical(model.ContentTypeJSON())
	response, statusCode := r.handlerService.Solds()
	if statusCode != http.StatusOK {
		json.NewEncoder(req).Encode(response)
		return
	}
	json.NewEncoder(req).Encode(response)
}

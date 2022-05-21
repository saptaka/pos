package api

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/go-playground/validator"
	"github.com/saptaka/pos/api/handler"
	"github.com/saptaka/pos/repository"
)

type Service interface {
	Route()
}

type service struct {
	routerHandler ApiRouter
}

func NewAPI(ctx context.Context, mux *router.Router, repo repository.Repo) Service {
	validation := validator.New()
	handlerService := handler.NewHandler(ctx, repo, validation)
	routerHandler := &apiRouter{handlerService, mux}
	return &service{routerHandler}
}

func (s *service) Route() {
	s.routerHandler.RouteCashierPath()
	s.routerHandler.RouteCategoryPath()
	s.routerHandler.RouteLoginPath()
	s.routerHandler.RoutePaymentPath()
	s.routerHandler.RouteProductPath()
	s.routerHandler.RouteReportPath()
	s.routerHandler.RouteOrderPath()
}

type apiRouter struct {
	handlerService handler.Service
	mux            *router.Router
}

type ApiRouter interface {
	CashierRouter
	CategoryRouter
	LoginRouter
	ProductRouter
	CategoryRouter
	PaymentRouter
	OrderRouter
	ReportRouter
}

func NewRouter() ApiRouter {
	return &apiRouter{}
}

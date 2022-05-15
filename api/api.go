package api

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/saptaka/pos/api/handler"
	"github.com/saptaka/pos/repository"
)

type Service interface {
	Route()
}

type service struct {
	routerHandler Router
}

func NewAPI(ctx context.Context, mux *mux.Router, repo repository.Repo) Service {
	validation := validator.New()
	handlerService := handler.NewHandler(ctx, repo, validation)
	routerHandler := &router{handlerService, mux}
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

type router struct {
	handlerService handler.Service
	mux            *mux.Router
}

type Router interface {
	CashierRouter
	CategoryRouter
	LoginRouter
	ProductRouter
	CategoryRouter
	PaymentRouter
	OrderRouter
	ReportRouter
}

func NewRouter() Router {
	return &router{}
}

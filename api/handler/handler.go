package handler

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/saptaka/pos/repository"
)

type Service interface {
	Cashier
	Login
	Category
	Product
	Payment
	Order
	Report
}

type service struct {
	ctx context.Context
	db  repository.Repo
}

var productCache syncMap

func NewHandler(ctx context.Context, db repository.Repo, validation *validator.Validate) Service {

	handlerService := service{ctx, db}
	productCache = syncMap{}
	go func() {
		err := handlerService.LoadProduct()
		if err != nil {
			panic(err)
		}
	}()
	return handlerService
}

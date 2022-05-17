package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type Product interface {
	ListProduct(limit, skip int, query string) ([]byte, int)
	DetailProduct(id int) ([]byte, int)
	CreateProduct(product model.Product) ([]byte, int)
	UpdateProduct(product model.Product) ([]byte, int)
	DeleteProduct(id int) ([]byte, int)
}

func (s service) ListProduct(limit, skip int, query string) ([]byte, int) {
	products, err := s.db.GetProducts(context.Background(), limit, skip, query)
	listProduct := model.ListProduct{
		Products: products,
		Meta: model.Meta{
			Total: len(products),
			Limit: limit,
			Skip:  skip,
		},
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, listProduct)
	}
	return utils.ResponseWrapper(http.StatusOK, listProduct)
}

func (s service) DetailProduct(id int) ([]byte, int) {
	Product, err := s.db.GetProductByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, Product)
}

func (s service) CreateProduct(productRequest model.Product) ([]byte, int) {
	err := s.validation.Struct(productRequest)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}

	product, err := s.db.CreateProduct(s.ctx, productRequest)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, product)
	}
	return utils.ResponseWrapper(http.StatusOK, product)
}

func (s service) UpdateProduct(product model.Product) ([]byte, int) {
	err := s.db.UpdateProduct(s.ctx, product)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

func (s service) DeleteProduct(id int) ([]byte, int) {
	err := s.db.DeleteProduct(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type syncMap struct {
	m sync.Map
}

func (c *syncMap) Get(key int64) (*model.Product, bool) {
	value, ok := c.m.Load(key)
	if ok {
		product := value.(*model.Product)
		return product, true
	}
	return nil, false
}

func (c *syncMap) Set(key int64, value *model.Product) {
	c.m.Store(key, value)
}

type Product interface {
	ListProduct(limit, skip int, categoryID int64, query string) ([]byte, int)
	DetailProduct(id int64) ([]byte, int)
	CreateProduct(product model.Product) ([]byte, int)
	UpdateProduct(product model.Product) ([]byte, int)
	DeleteProduct(id int64) ([]byte, int)
}

func (s service) ListProduct(limit, skip int, categoryID int64, query string) ([]byte, int) {
	products, err := s.db.GetProducts(context.Background(), limit, skip, categoryID, query)
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
		return utils.ResponseWrapper(http.StatusBadRequest, listProduct)
	}
	return utils.ResponseWrapper(http.StatusOK, listProduct)
}

func (s service) DetailProduct(id int64) ([]byte, int) {
	Product, err := s.db.GetProductByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
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

	productCache.Set(product.ProductId, &product)

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

	productCache.Set(product.ProductId, &product)

	return utils.ResponseWrapper(http.StatusOK, nil)
}

func (s service) DeleteProduct(id int64) ([]byte, int) {
	err := s.db.DeleteProduct(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

func (s service) LoadProduct() error {
	products, err := s.db.GetProducts(s.ctx, 0, 0, 0, "")
	for _, product := range products {
		productCache.Set(product.ProductId, &product)
	}
	return err
}

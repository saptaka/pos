package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

type syncMap struct {
	m sync.Map
}

func (c *syncMap) Get(key int64) (model.Product, bool) {
	value, ok := c.m.Load(key)
	if ok {
		product := value.(model.Product)
		return product, true
	}
	return model.Product{}, false
}

func (c *syncMap) Set(key int64, value model.Product) {
	c.m.Store(key, value)
}

type Product interface {
	ListProduct(limit, skip int, product model.Product) (map[string]interface{}, int)
	DetailProduct(id int64) (map[string]interface{}, int)
	CreateProduct(product model.ProductCreateRequest) (map[string]interface{}, int)
	UpdateProduct(product model.Product) (map[string]interface{}, int)
	DeleteProduct(id int64) (map[string]interface{}, int)
}

func (s service) ListProduct(limit, skip int, product model.Product) (map[string]interface{}, int) {
	products, err := s.db.GetProducts(s.ctx, limit, skip, product)

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
		return utils.ResponseWrapper(http.StatusBadRequest, listProduct, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, listProduct, nil)
}

func (s service) DetailProduct(id int64) (map[string]interface{}, int) {
	Product, err := s.db.GetProductByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, Product, nil)
}

func (s service) CreateProduct(productRequest model.ProductCreateRequest) (map[string]interface{}, int) {
	validation := productValidation(model.CREATE)
	err := validation.Struct(productRequest)
	if err != nil {
		return utils.ErrorWrapper(err, fasthttp.StatusBadRequest, model.CREATE)
	}
	product, err := s.db.CreateProduct(s.ctx, productRequest)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, product, nil)
	}

	productCache.Set(product.ProductId, product)

	productCreatedResponse := model.ProductCreateResponse{
		Name:       product.Name,
		ProductId:  product.ProductId,
		Stock:      product.Stock,
		SKU:        product.SKU,
		Price:      product.Price,
		Image:      product.Image,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
		CategoryId: product.CategoryId,
	}

	return utils.ResponseWrapper(http.StatusOK, productCreatedResponse, nil)
}

func (s service) UpdateProduct(product model.Product) (map[string]interface{}, int) {
	validation := productValidation(model.UPDATE)
	err := validation.Struct(product)
	if err != nil {
		return utils.ErrorWrapper(err, fasthttp.StatusBadRequest, model.UPDATE)
	}
	err = s.db.UpdateProduct(s.ctx, product)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}

	productCache.Set(product.ProductId, product)

	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

func (s service) DeleteProduct(id int64) (map[string]interface{}, int) {
	err := s.db.DeleteProduct(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

func (s service) LoadProduct() error {
	products, err := s.db.GetProducts(s.ctx, 0, 0, model.Product{})
	for _, product := range products {
		productCache.Set(product.ProductId, product)
	}
	return err
}

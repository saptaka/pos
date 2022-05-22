package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type Category interface {
	ListCategory(limit, skip int) (map[string]interface{}, int)
	DetailCategory(id int64) (map[string]interface{}, int)
	CreateCategory(Category model.Category) (map[string]interface{}, int)
	UpdateCategory(Category model.Category) (map[string]interface{}, int)
	DeleteCategory(id int64) (map[string]interface{}, int)
}

func (s service) ListCategory(limit, skip int) (map[string]interface{}, int) {
	categories, err := s.db.GetCategories(context.Background(), limit, skip)

	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	listCashier := model.ListCategory{
		Categories: categories,
		Meta: model.Meta{
			Total: len(categories),
			Limit: limit,
			Skip:  skip,
		},
	}
	return utils.ResponseWrapper(http.StatusOK, listCashier, nil)
}

func (s service) DetailCategory(id int64) (map[string]interface{}, int) {
	category, err := s.db.GetCategoryByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, category, nil)
}

func (s service) CreateCategory(category model.Category) (map[string]interface{}, int) {
	err := s.validation.Struct(category)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	CategoryData, err := s.db.CreateCategory(s.ctx, category.Name)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, CategoryData, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, CategoryData, nil)
}

func (s service) UpdateCategory(category model.Category) (map[string]interface{}, int) {

	err := s.db.UpdateCategory(s.ctx, category)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

func (s service) DeleteCategory(id int64) (map[string]interface{}, int) {
	err := s.db.DeleteCategory(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil, nil)
}

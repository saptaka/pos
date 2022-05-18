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
	ListCategory(limit, skip int) ([]byte, int)
	DetailCategory(id int64) ([]byte, int)
	CreateCategory(Category model.Category) ([]byte, int)
	UpdateCategory(Category model.Category) ([]byte, int)
	DeleteCategory(id int64) ([]byte, int)
}

func (s service) ListCategory(limit, skip int) ([]byte, int) {
	categories, err := s.db.GetCategories(context.Background(), limit, skip)

	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	listCashier := model.ListCategory{
		Categories: categories,
		Meta: model.Meta{
			Total: len(categories),
			Limit: limit,
			Skip:  skip,
		},
	}
	return utils.ResponseWrapper(http.StatusOK, listCashier)
}

func (s service) DetailCategory(id int64) ([]byte, int) {
	category, err := s.db.GetCategoryByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, category)
}

func (s service) CreateCategory(category model.Category) ([]byte, int) {
	err := s.validation.Struct(category)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	CategoryData, err := s.db.CreateCategory(s.ctx, category.Name)
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, CategoryData)
	}
	return utils.ResponseWrapper(http.StatusOK, CategoryData)
}

func (s service) UpdateCategory(category model.Category) ([]byte, int) {

	err := s.db.UpdateCategory(s.ctx, category)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

func (s service) DeleteCategory(id int64) ([]byte, int) {
	err := s.db.DeleteCategory(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

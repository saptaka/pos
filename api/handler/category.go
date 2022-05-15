package handler

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/saptaka/pos/model"
	"github.com/saptaka/pos/utils"
)

type Category interface {
	ListCategory(limit, skip int) ([]byte, int)
	DetailCategory(id int) ([]byte, int)
	CreateCategory(Category model.Category) ([]byte, int)
	UpdateCategory(Category model.Category) ([]byte, int)
	DeleteCategory(id int) ([]byte, int)
}

func (s service) ListCategory(limit, skip int) ([]byte, int) {
	categories, err := s.db.GetCategories(context.Background(), limit, skip)
	listCashier := model.ListCategory{
		Categories: categories,
		Meta: model.Meta{
			Total: len(categories),
			Limit: limit,
			Skip:  skip,
		},
	}
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, listCashier)
	}
	return utils.ResponseWrapper(http.StatusOK, listCashier)
}

func (s service) DetailCategory(id int) ([]byte, int) {
	category, err := s.db.GetCategoryByID(context.Background(), id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, category)
}

func (s service) CreateCategory(category model.Category) ([]byte, int) {
	err := s.validation.Struct(category)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	CategoryData, err := s.db.CreateCategory(s.ctx, category.Name)
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, CategoryData)
	}
	return utils.ResponseWrapper(http.StatusOK, CategoryData)
}

func (s service) UpdateCategory(category model.Category) ([]byte, int) {
	err := s.validation.Struct(category)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	err = s.db.UpdateCategory(s.ctx, category)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

func (s service) DeleteCategory(id int) ([]byte, int) {
	err := s.db.DeleteCategory(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, nil)
}

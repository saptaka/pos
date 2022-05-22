package handler

import (
	"net/http"

	"github.com/saptaka/pos/utils"
)

type Report interface {
	Revenue() (map[string]interface{}, int)
	Solds() (map[string]interface{}, int)
}

func (s service) Revenue() (map[string]interface{}, int) {
	revenue, err := s.db.GetRevenues(s.ctx)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, revenue, nil)
}

func (s service) Solds() (map[string]interface{}, int) {
	sold, err := s.db.GetSolds(s.ctx)
	if err != nil {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, sold, nil)
}

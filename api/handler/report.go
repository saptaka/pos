package handler

import (
	"net/http"

	"github.com/saptaka/pos/utils"
)

type Report interface {
	Revenue() ([]byte, int)
	Solds() ([]byte, int)
}

func (s service) Revenue() ([]byte, int) {
	revenue, err := s.db.GetRevenues(s.ctx)
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	return utils.ResponseWrapper(http.StatusOK, revenue)
}

func (s service) Solds() ([]byte, int) {
	return utils.ResponseWrapper(http.StatusAccepted, "haha")
}

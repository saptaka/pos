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
	return utils.ResponseWrapper(http.StatusAccepted, "haha")
}

func (s service) Solds() ([]byte, int) {
	return utils.ResponseWrapper(http.StatusAccepted, "haha")
}

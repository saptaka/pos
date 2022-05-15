package handler

import (
	"database/sql"
	"net/http"

	"github.com/saptaka/pos/utils"
)

type Login interface {
	GetPasscode(id int64) ([]byte, int)
	VerifyLogin(id int64, passcode, token string) ([]byte, int)
	VerifyLogout(id int64, passcode string) ([]byte, int)
}

func (s service) GetPasscode(id int64) ([]byte, int) {
	passcode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	dataPasscode := make(map[string]interface{})
	dataPasscode["passcode"] = passcode
	return utils.ResponseWrapper(http.StatusOK, dataPasscode)
}
func (s service) VerifyLogin(id int64, passcode, token string) ([]byte, int) {
	cashierPasscode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	if cashierPasscode != passcode {
		return utils.ResponseWrapper(http.StatusForbidden, nil)
	}
	tokenData := make(map[string]interface{})
	tokenData["token"] = token
	return utils.ResponseWrapper(http.StatusOK, tokenData)
}
func (s service) VerifyLogout(id int64, passcode string) ([]byte, int) {
	cashierPasscode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil)
	}
	if err != nil {
		return utils.ResponseWrapper(http.StatusInternalServerError, nil)
	}
	if cashierPasscode == passcode {
		return utils.ResponseWrapper(http.StatusOK, nil)
	}
	return utils.ResponseWrapper(http.StatusBadRequest, nil)
}

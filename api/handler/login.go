package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

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
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	dataPasscode := make(map[string]interface{})
	dataPasscode["passcode"] = passcode
	return utils.ResponseWrapper(http.StatusOK, dataPasscode)
}
func (s service) VerifyLogin(id int64, passcode, token string) ([]byte, int) {
	cashierPasscode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if cashierPasscode != passcode {
		return utils.ResponseWrapper(http.StatusUnauthorized, nil)
	}
	tokenData := make(map[string]interface{})
	tokenData["token"] = token
	return utils.ResponseWrapper(http.StatusOK, tokenData)
}
func (s service) VerifyLogout(id int64, passcode string) ([]byte, int) {
	cashierPasscode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil)
	}
	if strings.Compare(cashierPasscode, passcode) == 0 {
		return utils.ResponseWrapper(http.StatusOK, nil)
	}
	return utils.ResponseWrapper(http.StatusForbidden, nil)
}

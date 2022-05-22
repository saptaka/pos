package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/saptaka/pos/utils"
)

type Login interface {
	GetPasscode(id int64) (map[string]interface{}, int)
	VerifyLogin(id int64, passcode, token string) (map[string]interface{}, int)
	VerifyLogout(id int64, passcode string) (map[string]interface{}, int)
}

func (s service) GetPasscode(id int64) (map[string]interface{}, int) {
	passcode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusNotFound, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	dataPasscode := make(map[string]interface{})
	dataPasscode["passcode"] = passcode
	return utils.ResponseWrapper(http.StatusOK, dataPasscode, nil)
}
func (s service) VerifyLogin(id int64, passcode, token string) (map[string]interface{}, int) {
	cashierPasscode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if cashierPasscode != passcode {
		return utils.ResponseWrapper(http.StatusUnauthorized, nil, nil)
	}
	tokenData := make(map[string]interface{})
	tokenData["token"] = token
	return utils.ResponseWrapper(http.StatusOK, tokenData, nil)
}
func (s service) VerifyLogout(id int64, passcode string) (map[string]interface{}, int) {
	cashierPasscode, err := s.db.GetPasscodeById(s.ctx, id)
	if err == sql.ErrNoRows {
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if err != nil {
		log.Println(err)
		return utils.ResponseWrapper(http.StatusBadRequest, nil, nil)
	}
	if strings.Compare(cashierPasscode, passcode) == 0 {
		return utils.ResponseWrapper(http.StatusOK, nil, nil)
	}
	return utils.ResponseWrapper(http.StatusForbidden, nil, nil)
}

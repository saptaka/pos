package utils

import (
	"encoding/json"
	"net/http"

	"github.com/saptaka/pos/model"
)

const (
	Success = "Success"
	Error   = "Error"
)

func ResponseWrapper(statusCode int, data interface{}) ([]byte, int) {

	var status bool
	var message string
	if statusCode == http.StatusOK {
		status = true
		message = Success
	} else {
		status = false
		message = Error
	}
	var responseMessage map[string]interface{}
	if data == nil {
		responseMessage = make(map[string]interface{})
		responseMessage["success"] = status
		responseMessage["message"] = message
		jsonData, err := json.Marshal(responseMessage)
		if err != nil {
			return nil, http.StatusInternalServerError
		}
		return jsonData, statusCode

	}
	response := model.Response{
		Success: status,
		Message: message,
		Data:    data,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return jsonData, statusCode
}

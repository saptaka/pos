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
	if statusCode != http.StatusOK {
		status = false
		message = Error
		response := model.ErrorResponse{
			Response: model.Response{
				Success: status,
				Message: message,
			},
			Error: "error",
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			return nil, http.StatusInternalServerError
		}
		return jsonData, statusCode
	}

	response := model.SuccessResponse{
		Response: model.Response{
			Success: status,
			Message: message,
		},
		Data: data,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return jsonData, statusCode
}

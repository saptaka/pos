package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/saptaka/pos/model"
)

const (
	Success = "Success"
	Error   = "Error"
)

func ResponseWrapper(statusCode int, data interface{}) ([]byte, int) {

	if statusCode != http.StatusOK {

		response := model.ErrorResponse{
			Response: model.Response{
				Success: false,
				Message: Error,
			},
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			return nil, http.StatusBadRequest
		}
		return jsonData, statusCode
	}

	response := model.SuccessResponse{
		Response: model.Response{
			Success: true,
			Message: Success,
		},
		Data: data,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return nil, http.StatusBadRequest
	}

	return jsonData, statusCode
}

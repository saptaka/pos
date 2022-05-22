package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"github.com/saptaka/pos/model"
)

const (
	Success = "Success"
	Error   = "Error"
)

func ResponseWrapper(statusCode int, data interface{},
	errorData []model.ErrorData) (map[string]interface{}, int) {
	var response map[string]interface{}
	if statusCode != http.StatusOK {

		errorResponse := model.ErrorResponse{
			Response: model.Response{
				Success: false,
				Message: data.(string),
			},
			Error: []interface{}{errorData},
		}
		jsonData, err := json.Marshal(errorResponse)
		if err != nil {
			log.Println(err)
			return nil, http.StatusBadRequest
		}

		err = json.Unmarshal(jsonData, &response)
		if err != nil {
			log.Println(err)
			return nil, http.StatusBadRequest
		}
		return response, statusCode
	}

	successResponse := model.SuccessResponse{
		Response: model.Response{
			Success: true,
			Message: Success,
		},
		Data: data,
	}
	jsonData, err := json.Marshal(successResponse)
	if err != nil {
		log.Println(err)
		return nil, http.StatusBadRequest
	}

	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		log.Println(err)
		return nil, http.StatusBadRequest
	}
	return response, statusCode

}

func FormatCommas(num int) string {
	str := fmt.Sprintf("%d", num)
	re := regexp.MustCompile(`(\d+)(\d{3})`)
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}

func ErrorWrapper(err error, statusCode int, processType string) (map[string]interface{}, int) {
	var errorDatas []model.ErrorData
	errorMessage := "body ValidationError: "
	var errorContext interface{}

	validationErrors := err.(validator.ValidationErrors)
	for _, validationError := range validationErrors {
		field := validationError.Field()
		switch strings.ToUpper(processType) {
		case "CREATE":
			errorContext = model.CreateErrorContext{
				Key:   field,
				Label: field,
			}

		case "UPDATE":
			errorContext = model.CustomErrorContext{
				Label: field,
				Peers: []string{
					field,
				},
				PeersWithLabels: []string{
					field,
				},
				Value: make(map[string]interface{}),
			}

		}
		param := validationError.Param()
		errorData := model.ErrorData{
			Message: param,
			Path:    []string{field},
			Type:    validationError.Tag(),
			Context: errorContext,
		}

		errorMessage += param
		errorDatas = append(errorDatas, errorData)
	}

	return ResponseWrapper(statusCode, errorMessage, errorDatas)
}

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
)

func ResponseWrapper(statusCode int, data interface{},
	errorData []model.ErrorData) (map[string]interface{}, int) {
	var response map[string]interface{}
	var errorMessage interface{}
	if errorData != nil {
		errorMessage = errorData
	} else {
		errorMessage = make(map[string]interface{})
	}
	if statusCode != http.StatusOK {
		var message string
		switch v := data.(type) {
		case string:
			message = data.(string)
			fmt.Printf("String: %v", v)
		}
		errorResponse := model.ErrorResponse{
			Response: model.Response{
				Success: false,
				Message: message,
			},
			Error: errorMessage,
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
	validationErrors := err.(validator.ValidationErrors)
	switch strings.ToUpper(processType) {
	case model.CREATE:
		for _, validationError := range validationErrors {
			field := validationError.Field()
			param := validationError.Param()
			errorContext := model.CreateErrorContext{
				Key:   field,
				Label: field,
			}
			errorData := model.ErrorData{
				Message: param,
				Path:    []string{field},
				Type:    validationError.Tag(),
				Context: errorContext,
			}
			errorMessage += param
			errorDatas = append(errorDatas, errorData)
		}

	case model.UPDATE:
		for _, validationError := range validationErrors {

			field := validationError.Field()
			fields := strings.Split(field, ",")
			param := validationError.Param()
			errorContext := model.CustomErrorContext{
				Label:           "value",
				Peers:           fields,
				PeersWithLabels: fields,
				Value:           make(map[string]interface{}),
			}
			errorData := model.ErrorData{
				Message: param,
				Path:    []string{},
				Type:    validationError.Tag(),
				Context: errorContext,
			}
			errorMessage += param
			errorDatas = append(errorDatas, errorData)
		}
	case model.SUBORDER:
		for _, validationError := range validationErrors {

			param := validationError.Param()
			errorContext := model.CustomErrorContext{
				Label: "value",
				Value: make(map[string]interface{}),
			}
			errorData := model.ErrorData{
				Message: param,
				Path:    []string{},
				Type:    validationError.Tag(),
				Context: errorContext,
			}
			errorMessage += param
			errorDatas = append(errorDatas, errorData)
		}
	}

	return ResponseWrapper(statusCode, errorMessage, errorDatas)
}

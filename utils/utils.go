package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/saptaka/pos/model"
)

const (
	Success = "Success"
	Error   = "Error"
)

func ResponseWrapper(statusCode int, data interface{}) (map[string]interface{}, int) {
	var response map[string]interface{}
	if statusCode != http.StatusOK {

		errorResponse := model.ErrorResponse{
			Response: model.Response{
				Success: false,
				Message: "body ValidationError: \"value\" must be an array",
			},
			Error: []interface{}{data},
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

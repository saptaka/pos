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

func ResponseWrapper(statusCode int, data interface{}) ([]byte, int) {

	if statusCode != http.StatusOK {

		response := model.ErrorResponse{
			Response: model.Response{
				Success: false,
				Message: Error,
			},
			Error: make([]interface{}, 0),
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

func FormatCommas(num int) string {
	str := fmt.Sprintf("%d", num)
	re := regexp.MustCompile(`(\d+)(\d{3})`)
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}

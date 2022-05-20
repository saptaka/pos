package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/saptaka/pos/utils"
)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NDkzMTU5NzksInN1YiI6MX0.Eb-zFl9pVL7lmVjJCf74SqUIhfe3VXJ3_uJhTvm7iYc"

func middleware(next func(res http.ResponseWriter, req *http.Request)) func(res http.ResponseWriter, req *http.Request) {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		printRequestBody(req)
		reqToken := req.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "JWT ")
		if len(splitToken) < 2 {
			_, statusCode := utils.ResponseWrapper(http.StatusUnauthorized, nil)
			log.Println("unknown token")
			res.WriteHeader(statusCode)
			return
		}
		reqToken = splitToken[1]

		if reqToken == "" {
			_, statusCode := utils.ResponseWrapper(http.StatusUnauthorized, nil)
			log.Println("unknown token")
			res.WriteHeader(statusCode)
			return
		}

		if reqToken != Token {
			_, statusCode := utils.ResponseWrapper(http.StatusUnauthorized, nil)
			log.Println("unknown token")
			res.WriteHeader(statusCode)
			return
		}

		next(res, req)
	})
}

func printRequestBody(req *http.Request) {
	var bodyBytes []byte
	var err error

	if req.Body != nil {
		bodyBytes, err = ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer req.Body.Close()
	}

	if len(bodyBytes) > 0 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
			fmt.Printf("JSON parse error: %v", err)
			return
		}
		fmt.Println(prettyJSON.String())
	}
}

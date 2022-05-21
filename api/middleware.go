package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/saptaka/pos/utils"
	"github.com/valyala/fasthttp"
)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NDkzMTU5NzksInN1YiI6MX0.Eb-zFl9pVL7lmVjJCf74SqUIhfe3VXJ3_uJhTvm7iYc"

func middleware(next func(req *fasthttp.RequestCtx)) func(req *fasthttp.RequestCtx) {
	return (func(req *fasthttp.RequestCtx) {
		token := req.Request.Header.Peek("Authorization")
		reqToken := string(token)
		splitToken := strings.Split(reqToken, "JWT ")
		if len(splitToken) < 2 {
			response, statusCode := utils.ResponseWrapper(http.StatusUnauthorized, nil)
			req.Response.SetStatusCode(statusCode)
			json.NewEncoder(req).Encode(response)
			return
		}
		reqToken = splitToken[1]

		if reqToken == "" {
			response, statusCode := utils.ResponseWrapper(http.StatusUnauthorized, nil)
			log.Println("unknown token")
			req.Response.SetStatusCode(statusCode)
			json.NewEncoder(req).Encode(response)
			return
		}

		if reqToken != Token {
			response, statusCode := utils.ResponseWrapper(http.StatusUnauthorized, nil)
			log.Println("unknown token")
			req.Response.SetStatusCode(statusCode)
			json.NewEncoder(req).Encode(response)
			return
		}

		next(req)
	})
}

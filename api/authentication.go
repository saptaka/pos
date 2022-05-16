package api

import (
	"net/http"
)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NDkzMTU5NzksInN1YiI6MX0.Eb-zFl9pVL7lmVjJCf74SqUIhfe3VXJ3_uJhTvm7iYc"

func verifyToken(next func(res http.ResponseWriter, req *http.Request)) func(res http.ResponseWriter, req *http.Request) {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// reqToken := req.Header.Get("Authorization")
		// splitToken := strings.Split(reqToken, "Bearer ")
		// if len(splitToken) < 2 {
		// 	_, statusCode := utils.ResponseWrapper(http.StatusForbidden, nil)
		// 	res.WriteHeader(statusCode)
		// 	return
		// }
		// reqToken = splitToken[1]

		// if reqToken == "" {
		// 	_, statusCode := utils.ResponseWrapper(http.StatusForbidden, nil)
		// 	res.WriteHeader(statusCode)
		// 	return
		// }

		// if reqToken != Token {
		// 	_, statusCode := utils.ResponseWrapper(http.StatusForbidden, nil)
		// 	res.WriteHeader(statusCode)
		// }
		next(res, req)
	})
}

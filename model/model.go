package model

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Response
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Response
	Error []interface{} `json:"error"`
}

type Meta struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

const RFC3339MilliZ = "2006-01-02T15:04:05.000Z07:00"

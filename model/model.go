package model

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Meta struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

type UserContext string

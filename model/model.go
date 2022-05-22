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

type ErrorData struct {
	Message string      `json:"message"`
	Path    []string    `json:"path"`
	Type    string      `json:"type"`
	Context interface{} `json:"context"`
}

type CreateErrorContext struct {
	Label string      `json:"label"`
	Key   interface{} `json:"key"`
}

func ContentTypeJSON() ([]byte, []byte) {
	var (
		strContentType     = []byte("Content-Type")
		strApplicationJSON = []byte("application/json")
	)
	return strContentType, strApplicationJSON
}

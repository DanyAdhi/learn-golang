package utils

import (
	"encoding/json"
	"net/http"
)

type MetaStruct struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseSuccessStruct struct {
	Meta MetaStruct  `json:"meta"`
	Data interface{} `json:"data"`
}

func ResponseSuccess(w http.ResponseWriter, code int, message string, data interface{}) {
	if message == "" {
		message = "success"
	}

	if code < 100 || code >= 600 {
		code = 200
	}

	response := ResponseSuccessStruct{
		Meta: MetaStruct{
			Success: true,
			Code:    code,
			Message: message,
		},
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	if message == "" {
		message = "Internal server error"
	}

	if code < 100 || code >= 600 {
		code = 500
	}

	response := MetaStruct{
		Success: false,
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

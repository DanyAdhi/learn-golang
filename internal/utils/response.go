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

type ResponseSuccessPaginationStruct struct {
	Meta       MetaStruct  `json:"meta"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}
type Pagination struct {
	TotalRecords int  `json:"total_records"`
	LimitRecords int  `json:"limit_records"`
	TotalPages   int  `json:"total_pages"`
	CurrentPage  int  `json:"current_page"`
	PrevPage     *int `json:"prev_page"`
	NextPage     *int `json:"next_page"`
}

func ResponseSuccessPagination(w http.ResponseWriter, code int, message string, data interface{}, pagination Pagination) {
	if message == "" {
		message = "success"
	}

	if code < 100 || code >= 600 {
		code = 200
	}

	response := ResponseSuccessPaginationStruct{
		Meta: MetaStruct{
			Success: true,
			Code:    code,
			Message: message,
		},
		Data:       data,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func GeneratePagination(limit int, page int, totalRecords int) *Pagination {

	// Hitung total pages
	totalPages := (totalRecords + limit - 1) / limit

	// Hitung prev_page dan next_page
	var prevPage *int
	var nextPage *int
	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}
	if page < totalPages {
		next := page + 1
		nextPage = &next
	}

	pagination := Pagination{
		TotalRecords: totalRecords,
		LimitRecords: limit,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PrevPage:     prevPage,
		NextPage:     nextPage,
	}
	return &pagination
}

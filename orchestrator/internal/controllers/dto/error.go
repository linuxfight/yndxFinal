package dto

type ApiError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

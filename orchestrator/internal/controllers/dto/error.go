package dto

type ApiError struct {
	Message string `json:"message" example:"error: nil ptr dereference"`
	Code    int    `json:"code" example:"500"`
}

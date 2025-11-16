package model

type ErrorDetail struct {
	Code    string `json:"code" example:"NOT_FOUND"`
	Message string `json:"message" example:"resource not found"`
}

type ErrorResponse struct {
	Error *ErrorDetail `json:"error"`
}

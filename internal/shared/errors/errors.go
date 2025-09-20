package errors

type ApiError struct {
	Code    int
	Message string
}

func NewApiError(code int, message string) *ApiError {
	return &ApiError{Code: code, Message: message}
}

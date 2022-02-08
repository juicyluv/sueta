package apperror

import (
	"fmt"
	"net/http"
)

var (
	// ErrNotFound is used when server needs to response with 404 Not Found status code.
	ErrNotFound = NewAppError(
		http.StatusNotFound,
		"a given resource is not found",
		"please, double check your request structure",
	)
)

// AppError describes a structure of an error response in JSON format.
type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developerMessage,omitempty"`
	HttpCode         int    `json:"code,omitempty"`
}

// NewAppError returns a new AppError instance.
func NewAppError(code int, message, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Message:          message,
		DeveloperMessage: developerMessage,
		HttpCode:         code,
	}
}

// BadRequestError returns a new AppError instance
// with 400 Bad Request status code.
func BadRequestError(message, developerMessage string) *AppError {
	return NewAppError(http.StatusBadRequest, message, developerMessage)
}

// InternalError returns a new AppError instance
// with 500 Internal Server Error status code.
func InternalError(message, developerMessage string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, developerMessage)
}

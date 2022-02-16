package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNotFound is used when server needs to response with 404 Not Found status code.
	ErrNotFound = NewAppError(
		http.StatusNotFound,
		"requested resource is not found",
		"please, double check your request",
	)

	// ErrNoRows is used when no rows returned from storage.
	ErrNoRows = errors.New("no rows")

	// ErrValidationFailed is used when input validation failed.
	ErrValidationFailed = errors.New("input validation failed. please, provide valid values")
)

// AppError describes a structure of an error response in JSON format.
type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developerMessage,omitempty"`
	HttpCode         int    `json:"code,omitempty"`
} // @name ErrorResponse

// NewAppError returns a new AppError instance.
func NewAppError(code int, message, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Message:          message,
		DeveloperMessage: developerMessage,
		HttpCode:         code,
	}
}

// Error returns a string representation of an error.
func (ae *AppError) Error() string {
	return ae.Err.Error()
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

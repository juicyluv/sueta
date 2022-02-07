package apperror

import "fmt"

var ()

// AppError describes a structure of an error response in JSON format.
type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developerMessage,omitempty"`
}

// NewAppError returns a new AppError instance.
func NewAppError(message, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

package domain

import (
	"net/http"
)

type CustomError interface {
	Error() string
	StatusCode() int
}

type customError struct {
	message    string
	statusCode int
}

func (e *customError) Error() string {
	return e.message
}

func (e *customError) StatusCode() int {
	return e.statusCode
}

type ValidationError struct {
	customError
}

func NewValidationError(message string) error {
	return &ValidationError{
		customError: customError{
			message:    message,
			statusCode: http.StatusBadRequest,
		},
	}
}

type ConflictError struct {
	customError
}

func NewConflictError(message string) error {
	return &ConflictError{
		customError: customError{
			message:    message,
			statusCode: http.StatusConflict,
		},
	}
}

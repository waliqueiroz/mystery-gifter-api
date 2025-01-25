package domain

import (
	"net/http"

	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type CustomError interface {
	Error() string
	StatusCode() int
	Details() any
}

type customError struct {
	message    string
	statusCode int
	details    any
}

func (e *customError) Error() string {
	return e.message
}

func (e *customError) StatusCode() int {
	return e.statusCode
}

func (e *customError) Details() any {
	return e.details
}

type ValidationError struct {
	customError
}

func NewValidationError(errs validator.ValidationErrors) error {
	return &ValidationError{
		customError: customError{
			message:    "validation failed",
			statusCode: http.StatusBadRequest,
			details:    errs,
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

type ResourceNotFoundError struct {
	customError
}

func NewResourceNotFoundError(message string) error {
	return &ResourceNotFoundError{
		customError: customError{
			message:    message,
			statusCode: http.StatusNotFound,
		},
	}
}

type UnauthorizedError struct {
	customError
}

func NewUnauthorizedError(message string) error {
	return &ResourceNotFoundError{
		customError: customError{
			message:    message,
			statusCode: http.StatusUnauthorized,
		},
	}
}

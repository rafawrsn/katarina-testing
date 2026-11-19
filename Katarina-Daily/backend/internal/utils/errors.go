package utils

import (
	"errors"
	"net/http"
)

// AppError is a structured error carrying an HTTP status and a stable
// machine-readable code, so handlers never have to guess how to render
// a service-layer failure.
type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// AsAppError extracts an *AppError from an error chain.
func AsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

func NewAppError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func WrapAppError(status int, code, message string, err error) *AppError {
	return &AppError{Status: status, Code: code, Message: message, Err: err}
}

// Common, reusable errors.
var (
	ErrNotFound = func(resource string) *AppError {
		return NewAppError(http.StatusNotFound, resource+"_NOT_FOUND", resource+" not found")
	}
	ErrUnauthorized = NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
	ErrForbidden    = NewAppError(http.StatusForbidden, "FORBIDDEN", "you do not have access to this resource")
	ErrValidation   = func(msg string) *AppError {
		return NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", msg)
	}
	ErrConflict = func(code, msg string) *AppError {
		return NewAppError(http.StatusConflict, code, msg)
	}
	ErrInternal = func(err error) *AppError {
		return WrapAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred", err)
	}
)

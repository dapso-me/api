package common

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound        = errors.New("NOT_FOUND")
	ErrUnauthorized    = errors.New("UNAUTHORIZED")
	ErrForbidden       = errors.New("FORBIDDEN")
	ErrTooManyRequests = errors.New("TOO_MANY_REQUESTS")
	ErrInfrastructure  = errors.New("INFRASTRUCTURE_ERROR")
)

type BusinessError interface {
	Error() string
	Code() string
	Message() string
}

type AppError struct {
	code      string
	message   string
	fullError error
}

func NewAppError(code, message string) *AppError {
	return &AppError{
		code:      code,
		message:   message,
		fullError: errors.New(message),
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("app error: %s", e.fullError.Error())
}

func (e *AppError) Unwrap() error {
	return e.fullError
}

func (e *AppError) Code() string {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

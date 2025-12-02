package errs

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
	ErrUnexpected = errors.New("unexpected error")
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(msg string) error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: msg,
		Err:     ErrNotFound,
	}
}

func NewValidationError(msg string) error {
	return AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: msg,
		Err:     ErrValidation,
	}
}

func NewUnexpectedError(msg string) error {
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Err:     ErrUnexpected,
	}
}

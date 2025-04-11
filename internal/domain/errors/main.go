package errors

import (
	"fmt"
)

func New(message string) error {
	return fmt.Errorf("%s", message)
}

type ValidationError struct {
	message string
}

func NewValidationError(message string) error {
	return &ValidationError{message: message}
}

func (e *ValidationError) Error() string {
	return e.message
}

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string) error {
	return &NotFoundError{message: message}
}

func (e *NotFoundError) Error() string {
	return e.message
}

type OperationError struct {
	message string
}

func NewOperationError(message string) error {
	return &OperationError{message: message}
}

func (e *OperationError) Error() string {
	return e.message
}

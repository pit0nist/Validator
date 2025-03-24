package homework

import (
	"errors"
	"fmt"
)

var (
	ErrNotStruct                   = errors.New("wrong argument given, should be a struct")
	ErrInvalidValidatorSyntax      = errors.New("invalid validator syntax")
	ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
	ErrLenValidationFailed         = errors.New("len validation failed")
	ErrInValidationFailed          = errors.New("in validation failed")
	ErrMaxValidationFailed         = errors.New("max validation failed")
	ErrMinValidationFailed         = errors.New("min validation failed")
)

type ValidationError struct {
	field string
	err   error
}

func NewValidationError(err error, field string) error {
	return &ValidationError{
		field: field,
		err:   err,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.err)
}

func (e *ValidationError) Unwrap() error {
	return e.err
}

func Validate(_ any) error {
	// TODO implement me
	return nil
}

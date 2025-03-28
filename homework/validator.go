package homework

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func Validate(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var errs []error

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !field.IsExported() && field.Tag.Get("validate") != "" {
			errs = append(errs, NewValidationError(ErrValidateForUnexportedFields, field.Name))
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		validators := strings.Split(tag, ";")
		for _, validator := range validators {
			validator = strings.TrimSpace(validator)
			if validator == "" {
				continue
			}

			parts := strings.SplitN(validator, ":", 2)
			if len(parts) != 2 {
				errs = append(errs, NewValidationError(ErrInvalidValidatorSyntax, field.Name))
				continue
			}

			key, value := parts[0], parts[1]
			var err error

			switch key {
			case "len":
				err = validateLen(fieldValue, value, field.Name)
			case "in":
				err = validateIn(fieldValue, value, field.Name)
			case "min":
				err = validateMin(fieldValue, value, field.Name)
			case "max":
				err = validateMax(fieldValue, value, field.Name)
			default:
				err = NewValidationError(ErrInvalidValidatorSyntax, field.Name)
			}

			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func validateLen(field reflect.Value, param, fieldName string) error {
	if field.Kind() != reflect.String {
		return NewValidationError(fmt.Errorf("len requires string, got %s", field.Kind()), fieldName)
	}

	length, err := strconv.Atoi(param)
	if err != nil || length < 0 { // Проверка на отрицательные значения
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	if field.Len() != length {
		return NewValidationError(ErrLenValidationFailed, fieldName)
	}
	return nil
}

func validateIn(field reflect.Value, param, fieldName string) error {
	if param == "" {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	options := strings.Split(param, ",")
	if len(options) == 0 {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	switch field.Kind() {
	case reflect.String:
		val := field.String()
		for _, opt := range options {
			if strings.TrimSpace(opt) == val {
				return nil
			}
		}
		return NewValidationError(ErrInValidationFailed, fieldName)

	case reflect.Int:
		val := field.Int()
		for _, opt := range options {
			trimmedOpt := strings.TrimSpace(opt)
			num, err := strconv.ParseInt(trimmedOpt, 10, 64)
			if err != nil {
				continue
			}
			if val == num {
				return nil
			}
		}
		return NewValidationError(ErrInValidationFailed, fieldName)

	default:
		return NewValidationError(fmt.Errorf("in requires string or int, got %s", field.Kind()), fieldName)
	}
}

func validateMin(field reflect.Value, param, fieldName string) error {
	min, err := strconv.Atoi(param)
	if err != nil {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	switch field.Kind() {
	case reflect.String:
		if field.Len() < min {
			return NewValidationError(ErrMinValidationFailed, fieldName)
		}
	case reflect.Int:
		if field.Int() < int64(min) {
			return NewValidationError(ErrMinValidationFailed, fieldName)
		}
	default:
		return NewValidationError(fmt.Errorf("min requires string or int, got %s", field.Kind()), fieldName)
	}
	return nil
}

func validateMax(field reflect.Value, param, fieldName string) error {
	max, err := strconv.Atoi(param)
	if err != nil {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	switch field.Kind() {
	case reflect.String:
		if field.Len() > max {
			return NewValidationError(ErrMaxValidationFailed, fieldName)
		}
	case reflect.Int:
		if field.Int() > int64(max) {
			return NewValidationError(ErrMaxValidationFailed, fieldName)
		}
	default:
		return NewValidationError(fmt.Errorf("max requires string or int, got %s", field.Kind()), fieldName)
	}
	return nil
}

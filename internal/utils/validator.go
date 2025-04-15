package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	custom_errors "github.com/application-ellas/ella-backend/internal/domain/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func ParseValidatorErrorMessage(err error) error {
	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		return errors.New("invalid validation error")
	}

	errorsMessage := []string{}

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			validationTag := e.Tag()
			if validationTag == "required_if_match_format" || validationTag == "required_if_not_match_format" {
				validationTag = "required"
			}
			errorsMessage = append(errorsMessage, fmt.Sprintf("Field %s(%s) failed validation: %s", e.StructField(), e.Kind(), validationTag))
		}
	}

	return custom_errors.NewValidationError(strings.Join(errorsMessage, ", "))
}

func ValidateRequiredIfFieldNotMatchFormat(fl validator.FieldLevel) bool {
	return !ValidateRequiredIfFieldMatchFormat(fl)
}

func ValidateRequiredIfFieldMatchFormat(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), " ")
	field, format := params[0], params[1]

	switch format {
	case "uuid":
		value := fl.Parent().FieldByName(field).String()
		_, err := uuid.Parse(value)
		if err != nil {
			return true
		}
		return checkIfTypeIsEmpty(fl.Field())
	default:
		return false
	}
}

func checkIfTypeIsEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		value := field.String()
		return value != ""
	default:
		return false
	}
}

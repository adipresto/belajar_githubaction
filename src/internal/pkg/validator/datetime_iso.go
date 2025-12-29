package internal_validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateISODate(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

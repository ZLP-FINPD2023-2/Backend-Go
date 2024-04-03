package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func customValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("isNotFutureDate", isNotFutureDate)
	return v
}

func isNotFutureDate(fldLvl validator.FieldLevel) bool {
	dateToValidate, ok := fldLvl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	dateToValidate = dateToValidate.UTC()
	currentDate := time.Now().UTC()

	return dateToValidate.Before(currentDate) || dateToValidate.Equal(currentDate)
}

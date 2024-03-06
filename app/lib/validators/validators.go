package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func CustomValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("isNotFutureDate", IsNotFutureDate)
	return v
}

func IsNotFutureDate(fldLvl validator.FieldLevel) bool {
	dateToValidate, ok := fldLvl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	dateToValidate = dateToValidate.UTC()
	currentDate := time.Now().UTC()

	return dateToValidate.Before(currentDate) || dateToValidate.Equal(currentDate)
}

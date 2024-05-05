package validators

import (
	"finapp/constants"
	"finapp/models"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

func customValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("isNotFutureDate", isNotFutureDate)
	v.RegisterValidation("periodicity", periodicityValidation)
	return v
}

func isNotFutureDate(fldLvl validator.FieldLevel) bool {
	dateToValidateStr := fldLvl.Field().String()
	dateToValidate, err := time.Parse(constants.DateFormat, dateToValidateStr)
	if err != nil {
		log.Println(err)
		return false
	}
	dateToValidate = dateToValidate.UTC()
	currentDate := time.Now().UTC()
	return dateToValidate.Before(currentDate) || dateToValidate.Equal(currentDate)
}

func periodicityValidation(fldLvl validator.FieldLevel) bool {
	periodicity := models.Periodicity(fldLvl.Field().String())
	switch periodicity {
	case models.PeriodicityMonthly:
		return true
	case models.PeriodicityDaily:
		return true
	case models.PeriodicityYearly:
		return true
	}
	return false
}

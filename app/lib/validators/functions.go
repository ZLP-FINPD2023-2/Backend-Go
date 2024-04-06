package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func IsValid(model any) error {
	validate := customValidator()
	if err := validate.Struct(model); err != nil {
		return err
	}
	return nil
}

func ParseValidationErrors(err error) map[string]string {
	vErr, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	errs := make(map[string]string)

	for _, f := range vErr {
		currErr := f.ActualTag()
		if f.Param() != "" {
			currErr = fmt.Sprintf("%s=%s", currErr, f.Param())
		}
		errs[f.Field()] = currErr
	}

	return errs
}

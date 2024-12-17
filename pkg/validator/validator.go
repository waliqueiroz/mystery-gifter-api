package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var goValidator = validator.New()

func Validate(data interface{}) error {
	err := goValidator.Struct(data)
	if err != nil {
		validationErrors := validator.ValidationErrors{}
		if errors.As(err, &validationErrors) {
			return errors.Join(validationErrors)
		}

		return err
	}

	return nil
}

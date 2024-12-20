package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	goValidator = validator.New()
	trans       ut.Translator
)

func init() {
	// Configures the translator and the default translations
	en := en.New()
	uni := ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(goValidator, trans)

	// Configures the function to get the JSON tag name
	goValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrors []FieldError

func Validate(data any) ValidationErrors {
	err := goValidator.Struct(data)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			fieldErrors := make(ValidationErrors, 0, len(validationErrors))

			for _, e := range validationErrors {
				fieldErrors = append(fieldErrors, FieldError{
					Field: e.Field(),
					Error: e.Translate(trans),
				})
			}

			return fieldErrors
		}

		return ValidationErrors{{Field: "general", Error: err.Error()}}
	}

	return ValidationErrors{}
}

package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type TestStruct struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=18"`
	Password string `json:"password" validate:"required,min=8"`
}

type NestedTestStruct struct {
	User     TestStruct `json:"user" validate:"required"`
	Metadata struct {
		Version string   `json:"version" validate:"required"`
		Tags    []string `json:"tags" validate:"required,min=1"`
	} `json:"metadata" validate:"required"`
}

func TestValidate(t *testing.T) {
	t.Run("should return empty validation errors when struct is valid", func(t *testing.T) {
		// given
		validStruct := TestStruct{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      25,
			Password: "password123",
		}

		// when
		errors := validator.Validate(validStruct)

		// then
		assert.Empty(t, errors)
	})

	t.Run("should return validation errors when required fields are missing", func(t *testing.T) {
		// given
		invalidStruct := TestStruct{}

		// when
		errors := validator.Validate(invalidStruct)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 4)
		assert.Contains(t, errors, validator.FieldError{Field: "name", Error: "name is a required field"})
		assert.Contains(t, errors, validator.FieldError{Field: "email", Error: "email is a required field"})
		assert.Contains(t, errors, validator.FieldError{Field: "age", Error: "age is a required field"})
		assert.Contains(t, errors, validator.FieldError{Field: "password", Error: "password is a required field"})
	})

	t.Run("should return validation errors when email is invalid", func(t *testing.T) {
		// given
		invalidStruct := TestStruct{
			Name:     "John Doe",
			Email:    "invalid-email",
			Age:      25,
			Password: "password123",
		}

		// when
		errors := validator.Validate(invalidStruct)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "email", Error: "email must be a valid email address"})
	})

	t.Run("should return validation errors when age is less than minimum", func(t *testing.T) {
		// given
		invalidStruct := TestStruct{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      16,
			Password: "password123",
		}

		// when
		errors := validator.Validate(invalidStruct)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "age", Error: "age must be 18 or greater"})
	})

	t.Run("should return validation errors when password is too short", func(t *testing.T) {
		// given
		invalidStruct := TestStruct{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      25,
			Password: "123",
		}

		// when
		errors := validator.Validate(invalidStruct)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "password", Error: "password must be at least 8 characters in length"})
	})

	t.Run("should return general error when validation fails for non-struct type", func(t *testing.T) {
		// given
		invalidData := "not a struct"

		// when
		errors := validator.Validate(invalidData)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "general", Error: "validator: (nil string)"})
	})

	t.Run("should validate nested structs successfully", func(t *testing.T) {
		// given
		validNestedStruct := NestedTestStruct{
			User: TestStruct{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      25,
				Password: "password123",
			},
			Metadata: struct {
				Version string   `json:"version" validate:"required"`
				Tags    []string `json:"tags" validate:"required,min=1"`
			}{
				Version: "1.0.0",
				Tags:    []string{"test"},
			},
		}

		// when
		errors := validator.Validate(validNestedStruct)

		// then
		assert.Empty(t, errors)
	})

	t.Run("should return validation errors for nested structs", func(t *testing.T) {
		// given
		invalidNestedStruct := NestedTestStruct{
			User: TestStruct{
				Name:     "John Doe",
				Email:    "invalid-email",
				Age:      16,
				Password: "123",
			},
			Metadata: struct {
				Version string   `json:"version" validate:"required"`
				Tags    []string `json:"tags" validate:"required,min=1"`
			}{
				Version: "",
				Tags:    []string{},
			},
		}

		// when
		errors := validator.Validate(invalidNestedStruct)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 5)
		assert.Contains(t, errors, validator.FieldError{Field: "email", Error: "email must be a valid email address"})
		assert.Contains(t, errors, validator.FieldError{Field: "age", Error: "age must be 18 or greater"})
		assert.Contains(t, errors, validator.FieldError{Field: "password", Error: "password must be at least 8 characters in length"})
		assert.Contains(t, errors, validator.FieldError{Field: "version", Error: "version is a required field"})
		assert.Contains(t, errors, validator.FieldError{Field: "tags", Error: "tags must contain at least 1 item"})
	})

	t.Run("should handle nil pointer to struct", func(t *testing.T) {
		// given
		var nilStruct *TestStruct

		// when
		errors := validator.Validate(nilStruct)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "general", Error: "validator: (nil *validator_test.TestStruct)"})
	})

	t.Run("should handle empty slice", func(t *testing.T) {
		// given
		type SliceStruct struct {
			Items []string `json:"items" validate:"required,min=1"`
		}
		emptySliceStruct := SliceStruct{
			Items: []string{},
		}

		// when
		errors := validator.Validate(emptySliceStruct)

		// then
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "items", Error: "items must contain at least 1 item"})
	})
}

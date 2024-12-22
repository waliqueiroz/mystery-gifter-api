package entrypoint_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
	"github.com/waliqueiroz/mystery-gifter-api/test/helper"
)

func TestCustomErrorHandler(t *testing.T) {
	t.Run("should handle fiber.Error and return the correct WebError response", func(t *testing.T) {
		// given
		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get("/test", func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "unprocessable entity")
		})

		req := httptest.NewRequest(fiber.MethodGet, "/test", nil)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "unprocessable_entity", result.Code)
		assert.Equal(t, "unprocessable entity", result.Message)
		assert.Nil(t, result.Details)
	})

	t.Run("should handle domain.ValidationError and return the correct WebError response", func(t *testing.T) {
		// given
		validationErrors := validator.ValidationErrors{
			{
				Field: "field",
				Error: "error detail",
			},
		}
		validationErr := domain.NewValidationError(validationErrors)
		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get("/test", func(c fiber.Ctx) error {
			return validationErr
		})

		req := httptest.NewRequest(fiber.MethodGet, "/test", nil)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "bad_request", result.Code)
		assert.Equal(t, "validation failed", result.Message)

		details, ok := result.Details.([]any)
		assert.True(t, ok)
		assert.Len(t, details, 1)

		detail, ok := details[0].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "field", detail["field"])
		assert.Equal(t, "error detail", detail["error"])
	})

	t.Run("should handle domain.ResourceNotFoundError and return the correct WebError response", func(t *testing.T) {
		// given
		resourceNotFoundErr := domain.NewResourceNotFoundError("resource not found")
		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get("/test", func(c fiber.Ctx) error {
			return resourceNotFoundErr
		})

		req := httptest.NewRequest(fiber.MethodGet, "/test", nil)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "not_found", result.Code)
		assert.Equal(t, "resource not found", result.Message)
		assert.Nil(t, result.Details)
	})

	t.Run("should handle domain.ConflictError and return the correct WebError response", func(t *testing.T) {
		// given
		conflictErr := domain.NewConflictError("conflict error")
		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get("/test", func(c fiber.Ctx) error {
			return conflictErr
		})

		req := httptest.NewRequest(fiber.MethodGet, "/test", nil)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusConflict, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "conflict", result.Code)
		assert.Equal(t, "conflict error", result.Message)
		assert.Nil(t, result.Details)
	})

	t.Run("should handle unexpected errors and return the correct WebError response", func(t *testing.T) {
		// given
		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get("/test", func(c fiber.Ctx) error {
			return errors.New("unexpected error")
		})

		req := httptest.NewRequest(fiber.MethodGet, "/test", nil)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "internal_server_error", result.Code)
		assert.Equal(t, "unexpected error", result.Message)
		assert.Nil(t, result.Details)
	})
}

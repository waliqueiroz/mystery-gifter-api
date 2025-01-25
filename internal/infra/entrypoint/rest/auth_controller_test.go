package rest_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application/mock_application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest/build_rest"
	"github.com/waliqueiroz/mystery-gifter-api/test/helper"
	"go.uber.org/mock/gomock"
)

func Test_AuthController_Login(t *testing.T) {
	route := "/api/login"

	t.Run("should authenticate user successfully", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "some_password"
		user := build_domain.NewUserBuilder().WithEmail(email).Build()
		userDTO := build_rest.NewUserDTOBuilder().WithID(user.ID).WithEmail(user.Email).WithName(user.Name).WithSurname(user.Surname).WithCreatedAt(user.CreatedAt).WithUpdatedAt(user.UpdatedAt).Build()

		credentials := build_domain.NewCredentialsBuilder().WithEmail(email).WithPassword(password).Build()
		credentialsDTO := build_rest.NewCredentialsDTOBuilder().WithEmail(email).WithPassword(password).Build()

		authSession := build_domain.NewAuthSessionBuilder().WithUser(user).Build()
		expectedAuthSession := build_rest.NewAuthSessionDTOBuilder().WithUser(userDTO).WithAccessToken(authSession.AccessToken).WithTokenType(authSession.TokenType).WithExpiresIn(authSession.ExpiresIn).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthService := mock_application.NewMockAuthService(mockCtrl)
		mockedAuthService.EXPECT().Login(gomock.Any(), credentials).Return(&authSession, nil)

		authController := rest.NewAuthController(mockedAuthService)

		payload := helper.EncodeJSON(t, credentialsDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, authController.Login)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.AuthSessionDTO
		helper.DecodeJSON(t, response.Body, &result)
		assert.Equal(t, expectedAuthSession, result)
	})

	t.Run("should return bad_request when it fails to map auth session from domain", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "some_password"
		user := build_domain.NewUserBuilder().WithEmail(email).Build()

		credentials := build_domain.NewCredentialsBuilder().WithEmail(email).WithPassword(password).Build()
		credentialsDTO := build_rest.NewCredentialsDTOBuilder().WithEmail(email).WithPassword(password).Build()

		authSession := build_domain.NewAuthSessionBuilder().WithUser(user).WithAccessToken("").Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthService := mock_application.NewMockAuthService(mockCtrl)
		mockedAuthService.EXPECT().Login(gomock.Any(), credentials).Return(&authSession, nil)

		authController := rest.NewAuthController(mockedAuthService)

		payload := helper.EncodeJSON(t, credentialsDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, authController.Login)

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
		assert.Equal(t, "access_token", detail["field"])
		assert.Equal(t, "access_token is a required field", detail["error"])
	})

	t.Run("should return internal_server_error when service fails with unexpected error", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "some_password"

		credentials := build_domain.NewCredentialsBuilder().WithEmail(email).WithPassword(password).Build()
		credentialsDTO := build_rest.NewCredentialsDTOBuilder().WithEmail(email).WithPassword(password).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthService := mock_application.NewMockAuthService(mockCtrl)
		mockedAuthService.EXPECT().Login(gomock.Any(), credentials).Return(nil, assert.AnError)

		authController := rest.NewAuthController(mockedAuthService)

		payload := helper.EncodeJSON(t, credentialsDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, authController.Login)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "internal_server_error", result.Code)
		assert.Equal(t, assert.AnError.Error(), result.Message)
	})

	t.Run("should return bad_request when credentials are invalid", func(t *testing.T) {
		// given
		password := "some_password"

		credentialsDTO := build_rest.NewCredentialsDTOBuilder().WithEmail("").WithPassword(password).Build()

		authController := rest.NewAuthController(nil)

		payload := helper.EncodeJSON(t, credentialsDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, authController.Login)

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
		assert.Equal(t, "email", detail["field"])
		assert.Equal(t, "email is a required field", detail["error"])
	})

	t.Run("should return unprocessable_entity when payload is malformed", func(t *testing.T) {
		// given
		authController := rest.NewAuthController(nil)

		payload := helper.EncodeJSON(t, "invalid_payload")

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, authController.Login)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "unprocessable_entity", result.Code)
		assert.Equal(t, "Unprocessable Entity", result.Message)
	})
}

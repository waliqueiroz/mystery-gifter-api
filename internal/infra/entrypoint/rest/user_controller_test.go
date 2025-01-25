package rest_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application/mock_application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest/build_rest"
	"github.com/waliqueiroz/mystery-gifter-api/test/helper"
	"go.uber.org/mock/gomock"
)

func Test_UserController_Create(t *testing.T) {
	route := "/api/users"

	t.Run("should return status 201 and the user ID when the user is created successfully", func(t *testing.T) {
		// given
		createUserDTO := build_rest.NewCreateUserDTOBuilder().Build()

		userID := uuid.New().String()
		hashedPassword := "some-hashed-password"
		user := build_domain.NewUserBuilder().
			WithID(userID).
			WithName(createUserDTO.Name).
			WithSurname(createUserDTO.Surname).
			WithEmail(createUserDTO.Email).
			WithPassword(hashedPassword).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(createUserDTO.Password).Return(hashedPassword, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(userID, nil)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, createdUser domain.User) error {
			createdUser.CreatedAt = user.CreatedAt
			createdUser.UpdatedAt = user.UpdatedAt
			assert.Equal(t, user, createdUser)
			return nil
		})

		userController := rest.NewUserController(mockedUserService, mockedIdentityGenerator, mockedPasswordManager)

		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, response.StatusCode)

		var result fiber.Map
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, userID, result["id"])
	})

	t.Run("should return internal_server_error with an error message when fail to create user", func(t *testing.T) {
		// given
		createUserDTO := build_rest.NewCreateUserDTOBuilder().Build()

		userID := uuid.New().String()
		hashedPassword := "some-hashed-password"
		user := build_domain.NewUserBuilder().
			WithID(userID).
			WithName(createUserDTO.Name).
			WithSurname(createUserDTO.Surname).
			WithEmail(createUserDTO.Email).
			WithPassword(hashedPassword).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(createUserDTO.Password).Return(hashedPassword, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(userID, nil)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, createdUser domain.User) error {
			createdUser.CreatedAt = user.CreatedAt
			createdUser.UpdatedAt = user.UpdatedAt
			assert.Equal(t, user, createdUser)
			return assert.AnError
		})

		userController := rest.NewUserController(mockedUserService, mockedIdentityGenerator, mockedPasswordManager)

		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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

	t.Run("should return internal_server_error with an error message when fail to generate user ID", func(t *testing.T) {
		// given
		createUserDTO := build_rest.NewCreateUserDTOBuilder().Build()

		hashedPassword := "some-hashed-password"

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(createUserDTO.Password).Return(hashedPassword, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return("", assert.AnError)

		userController := rest.NewUserController(nil, mockedIdentityGenerator, mockedPasswordManager)

		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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

	t.Run("should return internal_server_error with an error message when fail to generate user ID", func(t *testing.T) {
		// given
		createUserDTO := build_rest.NewCreateUserDTOBuilder().Build()

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(createUserDTO.Password).Return("", assert.AnError)

		userController := rest.NewUserController(nil, nil, mockedPasswordManager)

		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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

	t.Run("should return bad_request with an error message email is invalid", func(t *testing.T) {
		// given
		createUserDTO := build_rest.NewCreateUserDTOBuilder().WithEmail("invalid_email").Build()

		userController := rest.NewUserController(nil, nil, nil)

		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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
		assert.Equal(t, "email must be a valid email address", detail["error"])
	})

	t.Run("should return bad_request with an error message password is not equal to password_confirm", func(t *testing.T) {
		// given
		createUserDTO := build_rest.NewCreateUserDTOBuilder().WithPassword("12345678").WithPasswordConfirm("1234567").Build()

		userController := rest.NewUserController(nil, nil, nil)

		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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
		assert.Equal(t, "password", detail["field"])
		assert.Equal(t, "password must be equal to PasswordConfirm", detail["error"])
	})

	t.Run("should return bad_request with an error message password less than 8 characteres", func(t *testing.T) {
		// given
		createUserDTO := build_rest.NewCreateUserDTOBuilder().WithPassword("1234567").WithPasswordConfirm("1234567").Build()

		userController := rest.NewUserController(nil, nil, nil)

		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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
		assert.Equal(t, "password", detail["field"])
		assert.Equal(t, "password must be at least 8 characters in length", detail["error"])
	})

	t.Run("should return unprocessable_entity with an error message when receive an invalid payload", func(t *testing.T) {
		// given
		userController := rest.NewUserController(nil, nil, nil)

		payload := helper.EncodeJSON(t, "invalid_payload")

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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

	t.Run("should return unprocessable_entity with an error message when receive an empty payload", func(t *testing.T) {
		// given
		userController := rest.NewUserController(nil, nil, nil)

		req := httptest.NewRequest(fiber.MethodPost, route, nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, userController.Create)

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

func Test_UserController_GetByID(t *testing.T) {
	route := "/api/users/:userID"

	t.Run("should return status 200 and the user when the user is found successfully", func(t *testing.T) {
		// given
		userID := uuid.New().String()
		user := build_domain.NewUserBuilder().WithID(userID).Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), userID).Return(&user, nil)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/users/%s", userID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.GetByID)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.UserDTO
		helper.DecodeJSON(t, response.Body, &result)

		expectedUser := build_rest.NewUserDTOBuilder().
			WithID(user.ID).
			WithName(user.Name).
			WithSurname(user.Surname).
			WithEmail(user.Email).
			WithCreatedAt(user.CreatedAt).
			WithUpdatedAt(user.UpdatedAt).Build()

		assert.Equal(t, expectedUser, result)
	})

	t.Run("should return bad_request with an error message when fails to map user from domain", func(t *testing.T) {
		// given
		userID := uuid.New().String()
		user := build_domain.NewUserBuilder().WithID(userID).WithName("").Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), userID).Return(&user, nil)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/users/%s", userID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.GetByID)

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
		assert.Equal(t, "name", detail["field"])
		assert.Equal(t, "name is a required field", detail["error"])
	})

	t.Run("should return internal_server_error with an error message when fail to get user", func(t *testing.T) {
		// given
		userID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), userID).Return(nil, assert.AnError)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/users/%s", userID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.GetByID)

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
}

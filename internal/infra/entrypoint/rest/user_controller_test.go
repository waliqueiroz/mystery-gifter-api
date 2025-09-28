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
	route := "/api/v1/users"

	t.Run("should return status 201 and the user when the user is created successfully", func(t *testing.T) {
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

		expectedUserDTO := build_rest.NewUserDTOBuilder().
			WithID(user.ID).
			WithName(user.Name).
			WithSurname(user.Surname).
			WithEmail(user.Email).
			WithCreatedAt(user.CreatedAt).
			WithUpdatedAt(user.UpdatedAt).Build()

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

		var result rest.UserDTO
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, expectedUserDTO.ID, result.ID)
		assert.Equal(t, expectedUserDTO.Name, result.Name)
		assert.Equal(t, expectedUserDTO.Surname, result.Surname)
		assert.Equal(t, expectedUserDTO.Email, result.Email)
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

	t.Run("should return internal_server_error with an error message when fail to hash password", func(t *testing.T) {
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
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "email",
			"error": "email must be a valid email address",
		})
	})

	t.Run("should return bad_request with an error message when password is not equal to password_confirm", func(t *testing.T) {
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
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "password",
			"error": "password must be equal to PasswordConfirm",
		})
	})

	t.Run("should return bad_request with an error when message password less than 8 characters", func(t *testing.T) {
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
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "password",
			"error": "password must be at least 8 characters in length",
		})
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
	route := "/api/v1/users/:userID"

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
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "name",
			"error": "name is a required field",
		})
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

func Test_UserController_Search(t *testing.T) {
	route := "/api/v1/users/search"

	t.Run("should return status 200 and search result when search is successful", func(t *testing.T) {
		// given
		users := []domain.User{
			build_domain.NewUserBuilder().WithName("John").WithSurname("Doe").Build(),
			build_domain.NewUserBuilder().WithName("Jane").WithSurname("Smith").Build(),
		}

		searchResult := build_domain.NewSearchResultBuilder[domain.User]().
			WithResult(users).
			WithTotal(2).
			WithLimit(15).
			WithOffset(0).
			Build()

		userFilters := build_domain.NewUserFiltersBuilder().WithName("John").Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().Search(gomock.Any(), userFilters).Return(&searchResult, nil)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, route+"?name=John", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		expectedUserDTOs := []rest.UserDTO{
			build_rest.NewUserDTOBuilder().
				WithID(users[0].ID).
				WithName(users[0].Name).
				WithSurname(users[0].Surname).
				WithEmail(users[0].Email).
				WithCreatedAt(users[0].CreatedAt).
				WithUpdatedAt(users[0].UpdatedAt).
				Build(),
			build_rest.NewUserDTOBuilder().
				WithID(users[1].ID).
				WithName(users[1].Name).
				WithSurname(users[1].Surname).
				WithEmail(users[1].Email).
				WithCreatedAt(users[1].CreatedAt).
				WithUpdatedAt(users[1].UpdatedAt).
				Build(),
		}

		expectedSearchResultDTO := build_rest.NewSearchResultDTOBuilder[rest.UserDTO]().
			WithResult(expectedUserDTOs).
			WithTotal(searchResult.Paging.Total).
			WithLimit(searchResult.Paging.Limit).
			WithOffset(searchResult.Paging.Offset).
			Build()

		var result rest.SearchResultDTO[rest.UserDTO]
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, expectedSearchResultDTO, result)
	})

	t.Run("should return status 200 and empty result when no users found", func(t *testing.T) {
		// given
		searchResult := build_domain.NewSearchResultBuilder[domain.User]().
			WithResult([]domain.User{}).
			WithTotal(0).
			WithLimit(15).
			WithOffset(0).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&searchResult, nil)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, route+"?name=NonExistent", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		expectedSearchResultDTO := build_rest.NewSearchResultDTOBuilder[rest.UserDTO]().
			WithResult([]rest.UserDTO{}).
			WithTotal(0).
			WithLimit(15).
			WithOffset(0).
			Build()

		var result rest.SearchResultDTO[rest.UserDTO]
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, expectedSearchResultDTO, result)
	})

	t.Run("should return status 200 with custom filters when valid query parameters are provided", func(t *testing.T) {
		// given
		users := []domain.User{
			build_domain.NewUserBuilder().WithName("Alice").WithSurname("Johnson").WithEmail("alice@example.com").Build(),
		}

		searchResult := build_domain.NewSearchResultBuilder[domain.User]().
			WithResult(users).
			WithTotal(1).
			WithLimit(10).
			WithOffset(5).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		expectedFilters := build_domain.NewUserFiltersBuilder().
			WithName("Alice").
			WithSurname("Johnson").
			WithEmail("alice@example.com").
			WithLimit(10).
			WithOffset(5).
			WithSortDirection(domain.SortDirectionTypeDesc).
			WithSortBy("name").
			Build()

		mockedUserService.EXPECT().Search(gomock.Any(), expectedFilters).Return(&searchResult, nil)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, route+"?name=Alice&surname=Johnson&email=alice@example.com&limit=10&offset=5&sort_direction=DESC&sort_by=name", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		expectedUserDTO := build_rest.NewUserDTOBuilder().
			WithID(users[0].ID).
			WithName(users[0].Name).
			WithSurname(users[0].Surname).
			WithEmail(users[0].Email).
			WithCreatedAt(users[0].CreatedAt).
			WithUpdatedAt(users[0].UpdatedAt).
			Build()

		expectedSearchResultDTO := build_rest.NewSearchResultDTOBuilder[rest.UserDTO]().
			WithResult([]rest.UserDTO{expectedUserDTO}).
			WithTotal(1).
			WithLimit(10).
			WithOffset(5).
			Build()

		var result rest.SearchResultDTO[rest.UserDTO]
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, expectedSearchResultDTO, result)
	})

	t.Run("should return bad_request with an error message when sort_direction is invalid", func(t *testing.T) {
		// given
		userController := rest.NewUserController(nil, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, route+"?sort_direction=INVALID", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "bad_request", result.Code)
		assert.Equal(t, "validation failed", result.Message)
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "sort_direction",
			"error": "sort_direction must be one of [ASC DESC]",
		})
	})

	t.Run("should return bad_request with an error message when sort_by is invalid", func(t *testing.T) {
		// given
		userController := rest.NewUserController(nil, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, route+"?sort_by=invalid_field", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "bad_request", result.Code)
		assert.Equal(t, "validation failed", result.Message)
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "sort_by",
			"error": "sort_by must be one of [name surname email created_at updated_at]",
		})
	})

	t.Run("should return internal_server_error with an error message when user service fails", func(t *testing.T) {
		// given
		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().Search(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, route+"?name=John", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

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

	t.Run("should return bad_request with an error message when fails to map search result from domain", func(t *testing.T) {
		// given
		users := []domain.User{
			build_domain.NewUserBuilder().WithName("").Build(), // Invalid user with empty name
		}

		searchResult := build_domain.NewSearchResultBuilder[domain.User]().
			WithResult(users).
			WithTotal(1).
			WithLimit(15).
			WithOffset(0).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&searchResult, nil)

		userController := rest.NewUserController(mockedUserService, nil, nil)

		req := httptest.NewRequest(fiber.MethodGet, route+"?name=John", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "bad_request", result.Code)
		assert.Equal(t, "validation failed", result.Message)
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "name",
			"error": "name is a required field",
		})
	})

	t.Run("should return unprocessable_entity with an error message when query parsing fails", func(t *testing.T) {
		// given
		userController := rest.NewUserController(nil, nil, nil)

		// Creating a request with invalid query parameter that will cause parsing to fail
		req := httptest.NewRequest(fiber.MethodGet, route+"?limit=invalid_number", nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, userController.Search)

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

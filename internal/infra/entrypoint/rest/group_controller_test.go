package rest_test

import (
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

func Test_GroupController_Create(t *testing.T) {
	route := "/api/groups"

	t.Run("should return status 201 and the group when created successfully", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		createGroupDTO := build_rest.NewCreateGroupDTOBuilder().Build()

		user := build_domain.NewUserBuilder().WithID(authUserID).Build()
		group := build_domain.NewGroupBuilder().WithName(createGroupDTO.Name).WithOwnerID(user.ID).WithUsers([]domain.User{user}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().Create(gomock.Any(), createGroupDTO.Name, authUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, createGroupDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.Create)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, response.StatusCode)

		var result rest.GroupDTO
		helper.DecodeJSON(t, response.Body, &result)

		expectedUserDTO := build_rest.NewUserDTOBuilder().
			WithID(user.ID).
			WithName(user.Name).
			WithEmail(user.Email).
			WithCreatedAt(user.CreatedAt).
			WithUpdatedAt(user.UpdatedAt).
			Build()

		expectedGroupDTO := build_rest.NewGroupDTOBuilder().
			WithID(group.ID).
			WithName(group.Name).
			WithUsers([]rest.UserDTO{expectedUserDTO}).
			WithOwnerID(group.OwnerID).
			WithStatus(string(group.Status)).
			WithCreatedAt(group.CreatedAt).
			WithUpdatedAt(group.UpdatedAt).
			Build()

		assert.Equal(t, expectedGroupDTO, result)
	})

	t.Run("should return unprocessable_entity when payload is malformed", func(t *testing.T) {
		// given
		groupController := rest.NewGroupController(nil, nil)

		payload := helper.EncodeJSON(t, "invalid_payload")

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.Create)

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

	t.Run("should return bad_request when createGroupDTO is invalid", func(t *testing.T) {
		// given
		createGroupDTO := build_rest.NewCreateGroupDTOBuilder().WithName("").Build()

		groupController := rest.NewGroupController(nil, nil)

		payload := helper.EncodeJSON(t, createGroupDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.Create)

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

	t.Run("should return internal_server_error when session manager fails", func(t *testing.T) {
		// given
		createGroupDTO := build_rest.NewCreateGroupDTOBuilder().Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return("", assert.AnError)

		groupController := rest.NewGroupController(nil, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, createGroupDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.Create)

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

	t.Run("should return internal_server_error when group service fails", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		createGroupDTO := build_rest.NewCreateGroupDTOBuilder().Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().Create(gomock.Any(), createGroupDTO.Name, authUserID).Return(nil, assert.AnError)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, createGroupDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.Create)

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

	t.Run("should return bad_request with an error message when fails to map group from domain", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		createGroupDTO := build_rest.NewCreateGroupDTOBuilder().Build()

		user := build_domain.NewUserBuilder().WithID(authUserID).Build()
		group := build_domain.NewGroupBuilder().WithName("").WithOwnerID(user.ID).WithUsers([]domain.User{user}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().Create(gomock.Any(), createGroupDTO.Name, authUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, createGroupDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.Create)

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
}

func Test_GroupController_GetByID(t *testing.T) {
	route := "/api/groups/:groupID"

	t.Run("should return status 200 and the group when found successfully", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		user := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithID(groupID).WithOwnerID(user.ID).WithUsers([]domain.User{user}).Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GetByID(gomock.Any(), groupID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, nil)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/groups/%s", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, groupController.GetByID)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.GroupDTO
		helper.DecodeJSON(t, response.Body, &result)

		expectedUserDTO := build_rest.NewUserDTOBuilder().
			WithID(user.ID).
			WithName(user.Name).
			WithEmail(user.Email).
			WithCreatedAt(user.CreatedAt).
			WithUpdatedAt(user.UpdatedAt).
			Build()

		expectedGroupDTO := build_rest.NewGroupDTOBuilder().
			WithID(group.ID).
			WithName(group.Name).
			WithUsers([]rest.UserDTO{expectedUserDTO}).
			WithOwnerID(group.OwnerID).
			WithStatus(string(group.Status)).
			WithCreatedAt(group.CreatedAt).
			WithUpdatedAt(group.UpdatedAt).
			Build()

		assert.Equal(t, expectedGroupDTO, result)
	})

	t.Run("should return bad_request with an error message when fails to map group from domain", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		user := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithName("").WithID(groupID).WithOwnerID(user.ID).WithUsers([]domain.User{user}).Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GetByID(gomock.Any(), groupID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, nil)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/groups/%s", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, groupController.GetByID)

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

	t.Run("should return internal_server_error when group service fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GetByID(gomock.Any(), groupID).Return(nil, assert.AnError)

		groupController := rest.NewGroupController(mockedGroupService, nil)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/groups/%s", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, groupController.GetByID)

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

func Test_GroupController_AddUser(t *testing.T) {
	route := "/api/groups/:groupID/users"

	t.Run("should return status 200 and the updated group when the user is added successfully", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()
		targetUserID := uuid.New().String()
		addUserDTO := build_rest.NewAddUserDTOBuilder().WithUserID(targetUserID).Build()

		user := build_domain.NewUserBuilder().WithID(targetUserID).Build()
		group := build_domain.NewGroupBuilder().WithID(groupID).WithUsers([]domain.User{user}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().AddUser(gomock.Any(), groupID, authUserID, targetUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, addUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/users", groupID), payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.AddUser)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.GroupDTO
		helper.DecodeJSON(t, response.Body, &result)

		expectedUserDTO := build_rest.NewUserDTOBuilder().
			WithID(user.ID).
			WithName(user.Name).
			WithEmail(user.Email).
			WithCreatedAt(user.CreatedAt).
			WithUpdatedAt(user.UpdatedAt).
			Build()

		expectedGroupDTO := build_rest.NewGroupDTOBuilder().
			WithID(group.ID).
			WithName(group.Name).
			WithUsers([]rest.UserDTO{expectedUserDTO}).
			WithOwnerID(group.OwnerID).
			WithStatus(string(group.Status)).
			WithCreatedAt(group.CreatedAt).
			WithUpdatedAt(group.UpdatedAt).
			Build()

		assert.Equal(t, expectedGroupDTO, result)
	})

	t.Run("should return unprocessable_entity when payload is malformed", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		groupController := rest.NewGroupController(nil, nil)

		payload := helper.EncodeJSON(t, "invalid_payload")

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/users", groupID), payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.AddUser)

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

	t.Run("should return bad_request when addUserDTO is invalid", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		addUserDTO := build_rest.NewAddUserDTOBuilder().WithUserID("invalid-uuid").Build()

		groupController := rest.NewGroupController(nil, nil)

		payload := helper.EncodeJSON(t, addUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/users", groupID), payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.AddUser)

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
			"field": "user_id",
			"error": "user_id must be a valid UUID",
		})
	})

	t.Run("should return internal_server_error when session manager fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		targetUserID := uuid.New().String()
		addUserDTO := build_rest.NewAddUserDTOBuilder().WithUserID(targetUserID).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return("", assert.AnError)

		groupController := rest.NewGroupController(nil, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, addUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/users", groupID), payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.AddUser)

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

	t.Run("should return internal_server_error when group service fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()
		targetUserID := uuid.New().String()
		addUserDTO := build_rest.NewAddUserDTOBuilder().WithUserID(targetUserID).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().AddUser(gomock.Any(), groupID, authUserID, targetUserID).Return(nil, assert.AnError)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, addUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/users", groupID), payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.AddUser)

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

	t.Run("should return bad_request with an error message when fails to map group from domain", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()
		targetUserID := uuid.New().String()
		addUserDTO := build_rest.NewAddUserDTOBuilder().WithUserID(targetUserID).Build()

		user := build_domain.NewUserBuilder().WithID(targetUserID).Build()
		group := build_domain.NewGroupBuilder().WithName("").WithID(groupID).WithUsers([]domain.User{user}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().AddUser(gomock.Any(), groupID, authUserID, targetUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		payload := helper.EncodeJSON(t, addUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/users", groupID), payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.AddUser)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var result entrypoint.WebError
		helper.DecodeJSON(t, response.Body, &result)

		assert.Equal(t, "bad_request", result.Code)
		assert.Len(t, result.Details, 1)
		assert.Contains(t, result.Details, map[string]any{
			"field": "name",
			"error": "name is a required field",
		})
	})
}

func Test_GroupController_RemoveUser(t *testing.T) {
	route := "/api/groups/:groupID/users/:userID"

	t.Run("should return status 200 and the updated group when the user is removed successfully", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()
		targetUserID := uuid.New().String()

		owner := build_domain.NewUserBuilder().WithID(authUserID).Build()
		group := build_domain.NewGroupBuilder().WithID(groupID).WithOwnerID(owner.ID).WithUsers([]domain.User{owner}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().RemoveUser(gomock.Any(), groupID, authUserID, targetUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/groups/%s/users/%s", groupID, targetUserID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Delete(route, groupController.RemoveUser)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.GroupDTO
		helper.DecodeJSON(t, response.Body, &result)

		expectedUserDTO := build_rest.NewUserDTOBuilder().
			WithID(owner.ID).
			WithName(owner.Name).
			WithEmail(owner.Email).
			WithCreatedAt(owner.CreatedAt).
			WithUpdatedAt(owner.UpdatedAt).
			Build()

		expectedGroupDTO := build_rest.NewGroupDTOBuilder().
			WithID(group.ID).
			WithName(group.Name).
			WithUsers([]rest.UserDTO{expectedUserDTO}).
			WithOwnerID(group.OwnerID).
			WithStatus(string(group.Status)).
			WithCreatedAt(group.CreatedAt).
			WithUpdatedAt(group.UpdatedAt).
			Build()

		assert.Equal(t, expectedGroupDTO, result)
	})

	t.Run("should return internal_server_error when session manager fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		targetUserID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return("", assert.AnError)

		groupController := rest.NewGroupController(nil, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/groups/%s/users/%s", groupID, targetUserID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Delete(route, groupController.RemoveUser)

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

	t.Run("should return internal_server_error when group service fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()
		targetUserID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().RemoveUser(gomock.Any(), groupID, authUserID, targetUserID).Return(nil, assert.AnError)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/groups/%s/users/%s", groupID, targetUserID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Delete(route, groupController.RemoveUser)

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

	t.Run("should return bad_request with an error message when fails to map group from domain", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()
		targetUserID := uuid.New().String()

		owner := build_domain.NewUserBuilder().WithID(authUserID).Build()
		group := build_domain.NewGroupBuilder().WithName("").WithID(groupID).WithOwnerID(owner.ID).WithUsers([]domain.User{owner}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().RemoveUser(gomock.Any(), groupID, authUserID, targetUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/api/groups/%s/users/%s", groupID, targetUserID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Delete(route, groupController.RemoveUser)

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
}

func Test_GroupController_GenerateMatches(t *testing.T) {
	route := "/api/groups/:groupID/matches"

	t.Run("should return status 200 and the group with generated matches when successfully", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()

		giver := build_domain.NewUserBuilder().WithID(authUserID).Build()
		receiver := build_domain.NewUserBuilder().WithID(uuid.New().String()).Build()
		match := build_domain.NewMatchBuilder().WithGiverID(giver.ID).WithReceiverID(receiver.ID).Build()

		group := build_domain.NewGroupBuilder().WithID(groupID).WithOwnerID(giver.ID).WithUsers([]domain.User{giver, receiver}).WithMatches([]domain.Match{match}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GenerateMatches(gomock.Any(), groupID, authUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/matches", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.GenerateMatches)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.GroupDTO
		helper.DecodeJSON(t, response.Body, &result)

		expectedGiverDTO := build_rest.NewUserDTOBuilder().
			WithID(giver.ID).
			WithName(giver.Name).
			WithEmail(giver.Email).
			WithCreatedAt(giver.CreatedAt).
			WithUpdatedAt(giver.UpdatedAt).
			Build()

		expectedReceiverDTO := build_rest.NewUserDTOBuilder().
			WithID(receiver.ID).
			WithName(receiver.Name).
			WithEmail(receiver.Email).
			WithCreatedAt(receiver.CreatedAt).
			WithUpdatedAt(receiver.UpdatedAt).
			Build()

		expectedMatchDTO := build_rest.NewMatchDTOBuilder().
			WithGiverID(match.GiverID).
			WithReceiverID(match.ReceiverID).
			Build()

		expectedGroupDTO := build_rest.NewGroupDTOBuilder().
			WithID(group.ID).
			WithName(group.Name).
			WithUsers([]rest.UserDTO{expectedGiverDTO, expectedReceiverDTO}).
			WithOwnerID(group.OwnerID).
			WithMatches([]rest.MatchDTO{expectedMatchDTO}).
			WithStatus(string(group.Status)).
			WithCreatedAt(group.CreatedAt).
			WithUpdatedAt(group.UpdatedAt).
			Build()

		assert.Equal(t, expectedGroupDTO, result)
	})

	t.Run("should return internal_server_error when session manager fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return("", assert.AnError)

		groupController := rest.NewGroupController(nil, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/matches", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.GenerateMatches)

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

	t.Run("should return internal_server_error when group service fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GenerateMatches(gomock.Any(), groupID, authUserID).Return(nil, assert.AnError)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/matches", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.GenerateMatches)

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

	t.Run("should return bad_request with an error message when fails to map group from domain", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()

		giver := build_domain.NewUserBuilder().WithID(authUserID).Build()
		receiver := build_domain.NewUserBuilder().Build()
		match := build_domain.NewMatchBuilder().WithGiverID(giver.ID).WithReceiverID(receiver.ID).Build()

		group := build_domain.NewGroupBuilder().WithName("").WithID(groupID).WithOwnerID(giver.ID).WithUsers([]domain.User{giver, receiver}).WithMatches([]domain.Match{match}).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GenerateMatches(gomock.Any(), groupID, authUserID).Return(&group, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/groups/%s/matches", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupController.GenerateMatches)

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
}

func Test_GroupController_GetUserMatch(t *testing.T) {
	route := "/api/groups/:groupID/matches/user"

	t.Run("should return status 200 and the user match when found successfully", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()

		userMatch := build_domain.NewUserBuilder().WithID(uuid.New().String()).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GetUserMatch(gomock.Any(), groupID, authUserID).Return(&userMatch, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/groups/%s/matches/user", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, groupController.GetUserMatch)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.UserDTO
		helper.DecodeJSON(t, response.Body, &result)

		expectedUserDTO := build_rest.NewUserDTOBuilder().
			WithID(userMatch.ID).
			WithName(userMatch.Name).
			WithEmail(userMatch.Email).
			WithCreatedAt(userMatch.CreatedAt).
			WithUpdatedAt(userMatch.UpdatedAt).
			Build()

		assert.Equal(t, expectedUserDTO, result)
	})

	t.Run("should return internal_server_error when auth token manager fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return("", assert.AnError)

		groupController := rest.NewGroupController(nil, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/groups/%s/matches/user", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, groupController.GetUserMatch)

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

	t.Run("should return internal_server_error when group service fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GetUserMatch(gomock.Any(), groupID, authUserID).Return(nil, assert.AnError)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/groups/%s/matches/user", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, groupController.GetUserMatch)

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

	t.Run("should return bad_request with an error message when fails to map user from domain", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		authUserID := uuid.New().String()

		userMatch := build_domain.NewUserBuilder().WithID(authUserID).WithName("").Build() // Invalid name for mapping to DTO

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupService := mock_application.NewMockGroupService(mockCtrl)
		mockedGroupService.EXPECT().GetUserMatch(gomock.Any(), groupID, authUserID).Return(&userMatch, nil)

		groupController := rest.NewGroupController(mockedGroupService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/api/groups/%s/matches/user", groupID), nil)

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Get(route, groupController.GetUserMatch)

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
}

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
	"github.com/waliqueiroz/mystery-gifter-api/test/helper"
	"go.uber.org/mock/gomock"
)

func Test_GroupInviteController_Create(t *testing.T) {
	route := "/api/v1/groups/:groupID/invites"

	t.Run("should return status 201 and the group invite when created successfully", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		groupID := uuid.New().String()
		groupInvite := build_domain.NewGroupInviteBuilder().WithGroupID(groupID).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupInviteService := mock_application.NewMockGroupInviteService(mockCtrl)
		mockedGroupInviteService.EXPECT().Create(gomock.Any(), groupID, authUserID).Return(&groupInvite, nil)

		groupInviteController := rest.NewGroupInviteController(mockedGroupInviteService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/v1/groups/%s/invites", groupID), nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupInviteController.Create)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, response.StatusCode)

		var result rest.GroupInviteDTO
		helper.DecodeJSON(t, response.Body, &result)
		assert.Equal(t, groupInvite.ID, result.ID)
		assert.Equal(t, groupInvite.GroupID, result.GroupID)
	})

	t.Run("should return status 403 when user is not group owner", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		groupID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupInviteService := mock_application.NewMockGroupInviteService(mockCtrl)
		mockedGroupInviteService.EXPECT().Create(gomock.Any(), groupID, authUserID).Return(nil, domain.NewForbiddenError("only the group owner can create invites"))

		groupInviteController := rest.NewGroupInviteController(mockedGroupInviteService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/v1/groups/%s/invites", groupID), nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupInviteController.Create)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, response.StatusCode)
	})

	t.Run("should return status 404 when group is not found", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		groupID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupInviteService := mock_application.NewMockGroupInviteService(mockCtrl)
		mockedGroupInviteService.EXPECT().Create(gomock.Any(), groupID, authUserID).Return(nil, domain.NewResourceNotFoundError("group not found"))

		groupInviteController := rest.NewGroupInviteController(mockedGroupInviteService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/v1/groups/%s/invites", groupID), nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupInviteController.Create)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
	})

	t.Run("should return status 409 when group is not open", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		groupID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupInviteService := mock_application.NewMockGroupInviteService(mockCtrl)
		mockedGroupInviteService.EXPECT().Create(gomock.Any(), groupID, authUserID).Return(nil, domain.NewConflictError("group is not open for invites"))

		groupInviteController := rest.NewGroupInviteController(mockedGroupInviteService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/v1/groups/%s/invites", groupID), nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupInviteController.Create)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusConflict, response.StatusCode)
	})
}

func Test_GroupInviteController_Join(t *testing.T) {
	route := "/api/v1/invites/:inviteID/join"

	t.Run("should return status 200 and the group when joined successfully", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		inviteID := uuid.New().String()
		owner := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).Build()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupInviteService := mock_application.NewMockGroupInviteService(mockCtrl)
		mockedGroupInviteService.EXPECT().JoinGroup(gomock.Any(), inviteID, authUserID).Return(&group, nil)

		groupInviteController := rest.NewGroupInviteController(mockedGroupInviteService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/v1/invites/%s/join", inviteID), nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupInviteController.Join)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var result rest.GroupDTO
		helper.DecodeJSON(t, response.Body, &result)
		assert.Equal(t, group.ID, result.ID)
	})

	t.Run("should return status 404 when invite is not found", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		inviteID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupInviteService := mock_application.NewMockGroupInviteService(mockCtrl)
		mockedGroupInviteService.EXPECT().JoinGroup(gomock.Any(), inviteID, authUserID).Return(nil, domain.NewResourceNotFoundError("group invite not found"))

		groupInviteController := rest.NewGroupInviteController(mockedGroupInviteService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/v1/invites/%s/join", inviteID), nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupInviteController.Join)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
	})

	t.Run("should return status 409 when invite is expired", func(t *testing.T) {
		// given
		authUserID := uuid.New().String()
		inviteID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedAuthTokenManager := mock_domain.NewMockAuthTokenManager(mockCtrl)
		mockedAuthTokenManager.EXPECT().GetAuthUserID(gomock.Any()).Return(authUserID, nil)

		mockedGroupInviteService := mock_application.NewMockGroupInviteService(mockCtrl)
		mockedGroupInviteService.EXPECT().JoinGroup(gomock.Any(), inviteID, authUserID).Return(nil, domain.NewConflictError("invite has expired"))

		groupInviteController := rest.NewGroupInviteController(mockedGroupInviteService, mockedAuthTokenManager)

		req := httptest.NewRequest(fiber.MethodPost, fmt.Sprintf("/api/v1/invites/%s/join", inviteID), nil)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New(fiber.Config{
			ErrorHandler: entrypoint.CustomErrorHandler,
		})
		app.Post(route, groupInviteController.Join)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusConflict, response.StatusCode)
	})
}

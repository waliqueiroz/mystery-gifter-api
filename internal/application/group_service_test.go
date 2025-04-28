package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application/mock_application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
	"go.uber.org/mock/gomock"
)

func Test_groupService_Create(t *testing.T) {
	t.Run("should create group successfully", func(t *testing.T) {
		// given
		name := "Test Group"
		owner := build_domain.NewUserBuilder().Build()
		ownerID := owner.ID
		expectedGroup := build_domain.NewGroupBuilder().WithName(name).WithOwnerID(ownerID).WithUsers([]domain.User{owner}).Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), ownerID).Return(&owner, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(expectedGroup.ID, nil)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, group domain.Group) error {
			group.CreatedAt = expectedGroup.CreatedAt
			group.UpdatedAt = expectedGroup.UpdatedAt

			assert.Equal(t, expectedGroup, group)

			return nil
		})

		groupService := application.NewGroupService(mockedGroupRepository, mockedUserService, mockedIdentityGenerator)

		// when
		result, err := groupService.Create(context.Background(), name, ownerID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedGroup.ID, result.ID)
		assert.Equal(t, expectedGroup.Name, result.Name)
		assert.Equal(t, expectedGroup.OwnerID, result.OwnerID)
		assert.Equal(t, expectedGroup.Users, result.Users)
	})

	t.Run("should return validation error when name is empty", func(t *testing.T) {
		// given
		name := ""
		owner := build_domain.NewUserBuilder().Build()
		ownerID := owner.ID
		expectedGroup := build_domain.NewGroupBuilder().WithName(name).WithOwnerID(ownerID).WithUsers([]domain.User{owner}).Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), ownerID).Return(&owner, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(expectedGroup.ID, nil)

		groupService := application.NewGroupService(nil, mockedUserService, mockedIdentityGenerator)

		// when
		result, err := groupService.Create(context.Background(), name, ownerID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ValidationError
		assert.ErrorAs(t, err, &expectedError)
		assert.Equal(t, "validation failed", expectedError.Error())
		errors := expectedError.Details()
		assert.Contains(t, errors, validator.FieldError{Field: "Name", Error: "Name is a required field"})
	})

	t.Run("should return error when fails to get owner", func(t *testing.T) {
		// given
		name := "Test Group"
		ownerID := "invalid-id"

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), ownerID).Return(nil, assert.AnError)

		groupService := application.NewGroupService(nil, mockedUserService, nil)

		// when
		result, err := groupService.Create(context.Background(), name, ownerID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when fails to generate ID", func(t *testing.T) {
		// given
		name := "Test Group"
		owner := build_domain.NewUserBuilder().Build()
		ownerID := owner.ID

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), ownerID).Return(&owner, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return("", assert.AnError)

		groupService := application.NewGroupService(nil, mockedUserService, mockedIdentityGenerator)

		// when
		result, err := groupService.Create(context.Background(), name, ownerID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when fails to create group", func(t *testing.T) {
		// given
		name := "Test Group"
		owner := build_domain.NewUserBuilder().Build()
		ownerID := owner.ID
		expectedGroup := build_domain.NewGroupBuilder().WithName(name).WithOwnerID(ownerID).WithUsers([]domain.User{owner}).Build()

		mockCtrl := gomock.NewController(t)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), ownerID).Return(&owner, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(expectedGroup.ID, nil)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, mockedUserService, mockedIdentityGenerator)

		// when
		result, err := groupService.Create(context.Background(), name, ownerID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_groupService_GetByID(t *testing.T) {
	t.Run("should get group successfully", func(t *testing.T) {
		// given
		expectedGroup := build_domain.NewGroupBuilder().Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), expectedGroup.ID).Return(&expectedGroup, nil)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GetByID(context.Background(), expectedGroup.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedGroup, result)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// given
		groupID := "some-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), groupID).Return(nil, assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GetByID(context.Background(), groupID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_groupService_AddUser(t *testing.T) {
	t.Run("should add user successfully", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		requesterID := group.OwnerID
		targetUser := build_domain.NewUserBuilder().Build()

		expectedGroup := build_domain.NewGroupBuilder().
			WithID(group.ID).
			WithName(group.Name).
			WithOwnerID(group.OwnerID).
			WithUsers(append(group.Users, targetUser)).
			WithCreatedAt(group.CreatedAt).
			WithUpdatedAt(group.UpdatedAt).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedGroup domain.Group) error {
			updatedGroup.UpdatedAt = expectedGroup.UpdatedAt
			assert.Equal(t, expectedGroup, updatedGroup)

			return nil
		})

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), targetUser.ID).Return(&targetUser, nil)

		groupService := application.NewGroupService(mockedGroupRepository, mockedUserService, nil)

		// when
		result, err := groupService.AddUser(context.Background(), group.ID, requesterID, targetUser.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedGroup.ID, result.ID)
		assert.Equal(t, expectedGroup.Name, result.Name)
		assert.Equal(t, expectedGroup.OwnerID, result.OwnerID)
		assert.Equal(t, expectedGroup.Users, result.Users)
	})

	t.Run("should return error when fails to get group", func(t *testing.T) {
		// given
		groupID := "group-id"
		requesterID := "requester-id"
		targetUserID := "target-user-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), groupID).Return(nil, assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.AddUser(context.Background(), groupID, requesterID, targetUserID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when fails to get target user", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		requesterID := group.OwnerID
		targetUserID := "target-user-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), targetUserID).Return(nil, assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, mockedUserService, nil)

		// when
		result, err := groupService.AddUser(context.Background(), group.ID, requesterID, targetUserID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when fails to update group", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		requesterID := group.OwnerID
		targetUser := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(assert.AnError)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), targetUser.ID).Return(&targetUser, nil)

		groupService := application.NewGroupService(mockedGroupRepository, mockedUserService, nil)

		// when
		result, err := groupService.AddUser(context.Background(), group.ID, requesterID, targetUser.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return forbidden error when requester is not the group owner and try to add other user", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		requesterID := "requester-id"
		targetUser := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), targetUser.ID).Return(&targetUser, nil)

		groupService := application.NewGroupService(mockedGroupRepository, mockedUserService, nil)

		// when
		result, err := groupService.AddUser(context.Background(), group.ID, requesterID, targetUser.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ForbiddenError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, err, "only the group owner can add other users")
	})
}

func Test_groupService_RemoveUser(t *testing.T) {
	t.Run("should remove user successfully", func(t *testing.T) {
		// given
		groupOwner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		initialGroup := build_domain.NewGroupBuilder().WithOwnerID(groupOwner.ID).WithUsers([]domain.User{groupOwner, targetUser}).Build()
		requesterID := groupOwner.ID

		expectedGroup := build_domain.NewGroupBuilder().
			WithID(initialGroup.ID).
			WithName(initialGroup.Name).
			WithOwnerID(initialGroup.OwnerID).
			WithUsers([]domain.User{groupOwner}).
			WithCreatedAt(initialGroup.CreatedAt).
			WithUpdatedAt(initialGroup.UpdatedAt).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), initialGroup.ID).Return(&initialGroup, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedGroup domain.Group) error {
			updatedGroup.UpdatedAt = expectedGroup.UpdatedAt
			assert.Equal(t, expectedGroup, updatedGroup)

			return nil
		})

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.RemoveUser(context.Background(), initialGroup.ID, requesterID, targetUser.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedGroup.ID, result.ID)
		assert.Equal(t, expectedGroup.Name, result.Name)
		assert.Equal(t, expectedGroup.OwnerID, result.OwnerID)
		assert.Equal(t, expectedGroup.Users, result.Users)
	})

	t.Run("should return error when fails to get group", func(t *testing.T) {
		// given
		groupID := "group-id"
		requesterID := "requester-id"
		targetUserID := "target-user-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), groupID).Return(nil, assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.RemoveUser(context.Background(), groupID, requesterID, targetUserID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when fails to update group", func(t *testing.T) {
		// given
		groupOwner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(groupOwner.ID).WithUsers([]domain.User{groupOwner, targetUser}).Build()
		requesterID := groupOwner.ID

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.RemoveUser(context.Background(), group.ID, requesterID, targetUser.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return forbidden error when requester is not the group owner and try to remove other user", func(t *testing.T) {
		// given
		groupOwner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(groupOwner.ID).WithUsers([]domain.User{groupOwner, targetUser}).Build()
		requesterID := "requester-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.RemoveUser(context.Background(), group.ID, requesterID, targetUser.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ForbiddenError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, err, "only the group owner can remove other users")
	})
}

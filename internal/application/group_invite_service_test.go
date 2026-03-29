package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application/mock_application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
	"go.uber.org/mock/gomock"
)

func Test_groupInviteService_Create(t *testing.T) {
	t.Run("should create group invite successfully", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithStatus(domain.GroupStatusOpen).Build()
		generatedID := uuid.New().String()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedGroupInviteRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(generatedID, nil)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.Create(context.Background(), group.ID, owner.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, generatedID, result.ID)
		assert.Equal(t, group.ID, result.GroupID)
		assert.WithinDuration(t, time.Now().Add(expiration), result.ExpiresAt, time.Second)
	})

	t.Run("should return not found error when group does not exist", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		requesterID := uuid.New().String()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), groupID).Return(nil, domain.NewResourceNotFoundError("group not found"))

		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedUserService := mock_application.NewMockUserService(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.Create(context.Background(), groupID, requesterID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var notFoundErr *domain.ResourceNotFoundError
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return forbidden error when requester is not the owner", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		requester := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithStatus(domain.GroupStatusOpen).Build()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedUserService := mock_application.NewMockUserService(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.Create(context.Background(), group.ID, requester.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var forbiddenErr *domain.ForbiddenError
		assert.ErrorAs(t, err, &forbiddenErr)
	})

	t.Run("should return conflict error when group is not open", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithStatus(domain.GroupStatusMatched).Build()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedUserService := mock_application.NewMockUserService(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.Create(context.Background(), group.ID, owner.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
	})
}

func Test_groupInviteService_JoinGroup(t *testing.T) {
	t.Run("should join group successfully", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		joiningUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithStatus(domain.GroupStatusOpen).Build()
		groupInvite := build_domain.NewGroupInviteBuilder().WithGroupID(group.ID).WithExpiresAt(time.Now().Add(1 * time.Hour)).Build()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedGroupInviteRepository.EXPECT().GetByID(gomock.Any(), groupInvite.ID).Return(&groupInvite, nil)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), joiningUser.ID).Return(&joiningUser, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.JoinGroup(context.Background(), groupInvite.ID, joiningUser.ID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result.Users, joiningUser)
	})

	t.Run("should return not found error when invite does not exist", func(t *testing.T) {
		// given
		inviteID := uuid.New().String()
		userID := uuid.New().String()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedGroupInviteRepository.EXPECT().GetByID(gomock.Any(), inviteID).Return(nil, domain.NewResourceNotFoundError("group invite not found"))

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.JoinGroup(context.Background(), inviteID, userID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var notFoundErr *domain.ResourceNotFoundError
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return conflict error when invite is expired", func(t *testing.T) {
		// given
		joiningUser := build_domain.NewUserBuilder().Build()
		groupInvite := build_domain.NewGroupInviteBuilder().WithExpiresAt(time.Now().Add(-1 * time.Hour)).Build()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedGroupInviteRepository.EXPECT().GetByID(gomock.Any(), groupInvite.ID).Return(&groupInvite, nil)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.JoinGroup(context.Background(), groupInvite.ID, joiningUser.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
		assert.EqualError(t, conflictErr, "invite has expired")
	})

	t.Run("should return conflict error when group is not open", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		joiningUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithStatus(domain.GroupStatusMatched).Build()
		groupInvite := build_domain.NewGroupInviteBuilder().WithGroupID(group.ID).WithExpiresAt(time.Now().Add(1 * time.Hour)).Build()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedGroupInviteRepository := mock_domain.NewMockGroupInviteRepository(mockCtrl)
		mockedGroupInviteRepository.EXPECT().GetByID(gomock.Any(), groupInvite.ID).Return(&groupInvite, nil)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().GetByID(gomock.Any(), joiningUser.ID).Return(&joiningUser, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)

		groupInviteService := application.NewGroupInviteService(mockedGroupInviteRepository, mockedGroupRepository, mockedUserService, mockedIdentityGenerator, expiration)

		// when
		result, err := groupInviteService.JoinGroup(context.Background(), groupInvite.ID, joiningUser.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
	})
}

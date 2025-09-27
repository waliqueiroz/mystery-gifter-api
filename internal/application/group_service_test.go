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

func Test_groupService_GenerateMatches(t *testing.T) {
	t.Run("should generate matches successfully for an even number of users", func(t *testing.T) {
		// given
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		user3 := build_domain.NewUserBuilder().Build()
		user4 := build_domain.NewUserBuilder().Build()
		initialGroup := build_domain.NewGroupBuilder().
			WithOwnerID(user1.ID).
			WithUsers([]domain.User{user1, user2, user3, user4}).
			Build()
		requesterID := user1.ID

		expectedGroup := build_domain.NewGroupBuilder().
			WithID(initialGroup.ID).
			WithName(initialGroup.Name).
			WithOwnerID(initialGroup.OwnerID).
			WithUsers(initialGroup.Users).
			WithMatches([]domain.Match{
				{GiverID: user1.ID, ReceiverID: user2.ID},
				{GiverID: user2.ID, ReceiverID: user3.ID},
				{GiverID: user3.ID, ReceiverID: user4.ID},
				{GiverID: user4.ID, ReceiverID: user1.ID},
			}).
			WithCreatedAt(initialGroup.CreatedAt).
			WithUpdatedAt(initialGroup.UpdatedAt).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), initialGroup.ID).Return(&initialGroup, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedGroup domain.Group) error {
			assert.Equal(t, expectedGroup.ID, updatedGroup.ID)
			assert.ElementsMatch(t, expectedGroup.Users, updatedGroup.Users)
			// Because GenerateMatches shuffles the users, we can't assert the exact order of matches.
			// We can only assert that the count is correct and that the matches themselves are valid.
			assert.Len(t, updatedGroup.Matches, 4)
			return nil
		})

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GenerateMatches(context.Background(), initialGroup.ID, requesterID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Matches, 4)
		assert.Equal(t, expectedGroup.ID, result.ID)

		// Check if the generated matches are valid (each user is a giver and a receiver once)
		gifters := make(map[string]bool)
		receivers := make(map[string]bool)
		allUsersMap := make(map[string]bool)
		for _, u := range initialGroup.Users {
			allUsersMap[u.ID] = true
		}

		for _, match := range result.Matches {
			assert.Contains(t, allUsersMap, match.GiverID)
			assert.Contains(t, allUsersMap, match.ReceiverID)
			assert.NotEqual(t, match.GiverID, match.ReceiverID)
			gifters[match.GiverID] = true
			receivers[match.ReceiverID] = true
		}
		assert.Len(t, gifters, 4)
		assert.Len(t, receivers, 4)
	})

	t.Run("should generate matches successfully for an odd number of users", func(t *testing.T) {
		// given
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		user3 := build_domain.NewUserBuilder().Build()
		initialGroup := build_domain.NewGroupBuilder().
			WithOwnerID(user1.ID).
			WithUsers([]domain.User{user1, user2, user3}).
			Build()
		requesterID := user1.ID

		expectedGroup := build_domain.NewGroupBuilder().
			WithID(initialGroup.ID).
			WithName(initialGroup.Name).
			WithOwnerID(initialGroup.OwnerID).
			WithUsers(initialGroup.Users).
			WithMatches([]domain.Match{
				{GiverID: user1.ID, ReceiverID: user2.ID},
				{GiverID: user2.ID, ReceiverID: user3.ID},
				{GiverID: user3.ID, ReceiverID: user1.ID},
			}).
			WithCreatedAt(initialGroup.CreatedAt).
			WithUpdatedAt(initialGroup.UpdatedAt).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), initialGroup.ID).Return(&initialGroup, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedGroup domain.Group) error {
			assert.Equal(t, expectedGroup.ID, updatedGroup.ID)
			assert.ElementsMatch(t, expectedGroup.Users, updatedGroup.Users)
			assert.Len(t, updatedGroup.Matches, 3)
			return nil
		})

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GenerateMatches(context.Background(), initialGroup.ID, requesterID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Matches, 3)
		assert.Equal(t, expectedGroup.ID, result.ID)

		// Check if the generated matches are valid (each user is a giver and a receiver once)
		gifters := make(map[string]bool)
		receivers := make(map[string]bool)
		allUsersMap := make(map[string]bool)
		for _, u := range initialGroup.Users {
			allUsersMap[u.ID] = true
		}

		for _, match := range result.Matches {
			assert.Contains(t, allUsersMap, match.GiverID)
			assert.Contains(t, allUsersMap, match.ReceiverID)
			assert.NotEqual(t, match.GiverID, match.ReceiverID)
			gifters[match.GiverID] = true
			receivers[match.ReceiverID] = true
		}
		assert.Len(t, gifters, 3)
		assert.Len(t, receivers, 3)
	})

	t.Run("should return error when fails to get group", func(t *testing.T) {
		// given
		groupID := "some-group-id"
		requesterID := "some-requester-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), groupID).Return(nil, assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GenerateMatches(context.Background(), groupID, requesterID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return conflict error when domain group fails to generate matches (not enough users)", func(t *testing.T) {
		// given
		groupOwner := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		initialGroup := build_domain.NewGroupBuilder().
			WithOwnerID(groupOwner.ID).
			WithUsers([]domain.User{groupOwner, user2}). // Only two users, will cause an error in domain.Group.GenerateMatches
			Build()
		requesterID := groupOwner.ID

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), initialGroup.ID).Return(&initialGroup, nil)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GenerateMatches(context.Background(), initialGroup.ID, requesterID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ConflictError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "group must have at least 3 users to generate matches")
	})

	t.Run("should return forbidden error when requester is not the group owner", func(t *testing.T) {
		// given
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		user3 := build_domain.NewUserBuilder().Build()
		initialGroup := build_domain.NewGroupBuilder().
			WithOwnerID(user1.ID).
			WithUsers([]domain.User{user1, user2, user3}).
			Build()
		requesterID := "not-owner-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), initialGroup.ID).Return(&initialGroup, nil)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GenerateMatches(context.Background(), initialGroup.ID, requesterID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ForbiddenError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "only the group owner can generate matches")
	})

	t.Run("should return error when fails to update group after generating matches", func(t *testing.T) {
		// given
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		user3 := build_domain.NewUserBuilder().Build()
		initialGroup := build_domain.NewGroupBuilder().
			WithOwnerID(user1.ID).
			WithUsers([]domain.User{user1, user2, user3}).
			Build()
		requesterID := user1.ID

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), initialGroup.ID).Return(&initialGroup, nil)
		mockedGroupRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GenerateMatches(context.Background(), initialGroup.ID, requesterID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_groupService_GetUserMatch(t *testing.T) {
	t.Run("should return user match successfully", func(t *testing.T) {
		// given
		requester := build_domain.NewUserBuilder().Build()
		matchedUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithUsers([]domain.User{requester, matchedUser}).
			WithMatches([]domain.Match{{GiverID: requester.ID, ReceiverID: matchedUser.ID}}).
			Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GetUserMatch(context.Background(), group.ID, requester.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &matchedUser, result)
	})

	t.Run("should return error when fails to get group", func(t *testing.T) {
		// given
		groupID := "some-group-id"
		requesterID := "some-requester-id"

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), groupID).Return(nil, assert.AnError)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GetUserMatch(context.Background(), groupID, requesterID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when domain group fails to get user match", func(t *testing.T) {
		// given
		requester := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithUsers([]domain.User{requester}). // No match for requester
			Build()

		mockCtrl := gomock.NewController(t)

		mockedGroupRepository := mock_domain.NewMockGroupRepository(mockCtrl)
		mockedGroupRepository.EXPECT().GetByID(gomock.Any(), group.ID).Return(&group, nil)

		groupService := application.NewGroupService(mockedGroupRepository, nil, nil)

		// when
		result, err := groupService.GetUserMatch(context.Background(), group.ID, requester.ID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ConflictError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "match not found for the given user")
	})
}

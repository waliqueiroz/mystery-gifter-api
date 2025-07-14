package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
	"go.uber.org/mock/gomock"
)

func Test_NewGroup(t *testing.T) {
	t.Run("should create a new group successfully", func(t *testing.T) {
		// given
		name := "Test Group"
		generatedID := uuid.New().String()
		owner := build_domain.NewUserBuilder().Build()
		now := time.Now()

		mockCtrl := gomock.NewController(t)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(generatedID, nil)

		// when
		group, err := domain.NewGroup(mockedIdentityGenerator, name, owner)

		// then
		assert.NoError(t, err)
		assert.Equal(t, generatedID, group.ID)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, owner.ID, group.OwnerID)
		assert.Equal(t, []domain.User{owner}, group.Users)
		assert.WithinDuration(t, now, group.CreatedAt, time.Second)
		assert.WithinDuration(t, now, group.UpdatedAt, time.Second)
	})

	t.Run("should return error when identity generator fails", func(t *testing.T) {
		// given
		name := "Test Group"
		owner := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return("", assert.AnError)

		// when
		group, err := domain.NewGroup(mockedIdentityGenerator, name, owner)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, group)
	})

	t.Run("should return validation error when name is empty", func(t *testing.T) {
		// given
		name := ""
		owner := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(uuid.New().String(), nil)

		// when
		group, err := domain.NewGroup(mockedIdentityGenerator, name, owner)

		// then
		assert.Nil(t, group)
		assert.Error(t, err)
		var validationErr *domain.ValidationError
		assert.ErrorAs(t, err, &validationErr)
		errors := validationErr.Details()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "Name", Error: "Name is a required field"})
	})
}

func Test_Group_AddUser(t *testing.T) {
	t.Run("should add user successfully when requester is owner", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).Build()
		originalUpdateTime := group.UpdatedAt

		// when
		err := group.AddUser(owner.ID, targetUser)

		// then
		assert.NoError(t, err)
		assert.Contains(t, group.Users, targetUser)
		assert.NotEqual(t, originalUpdateTime, group.UpdatedAt)
	})

	t.Run("should add user successfully when requester is the target user", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).Build()
		originalUpdateTime := group.UpdatedAt

		// when
		err := group.AddUser(targetUser.ID, targetUser)

		// then
		assert.NoError(t, err)
		assert.Contains(t, group.Users, targetUser)
		assert.NotEqual(t, originalUpdateTime, group.UpdatedAt)
	})

	t.Run("should not add duplicate user", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, targetUser}).
			Build()
		originalUpdateTime := group.UpdatedAt

		// when
		err := group.AddUser(owner.ID, targetUser)

		// then
		assert.NoError(t, err)
		assert.Len(t, group.Users, 2)
		assert.Equal(t, originalUpdateTime, group.UpdatedAt)
	})

	t.Run("should return forbidden error when requester is not owner or target user", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		requester := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).Build()

		// when
		err := group.AddUser(requester.ID, targetUser)

		// then
		assert.Error(t, err)
		var forbiddenErr *domain.ForbiddenError
		assert.ErrorAs(t, err, &forbiddenErr)
		assert.EqualError(t, forbiddenErr, "only the group owner can add other users")
		assert.NotContains(t, group.Users, targetUser)
	})
}

func Test_Group_RemoveUser(t *testing.T) {
	t.Run("should remove user successfully when requester is owner", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, targetUser}).
			Build()
		originalUpdateTime := group.UpdatedAt

		// when
		err := group.RemoveUser(owner.ID, targetUser.ID)

		// then
		assert.NoError(t, err)
		assert.NotContains(t, group.Users, targetUser)
		assert.NotEqual(t, originalUpdateTime, group.UpdatedAt)
	})

	t.Run("should remove user successfully when requester is the target user", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, targetUser}).
			Build()
		originalUpdateTime := group.UpdatedAt

		// when
		err := group.RemoveUser(targetUser.ID, targetUser.ID)

		// then
		assert.NoError(t, err)
		assert.NotContains(t, group.Users, targetUser)
		assert.NotEqual(t, originalUpdateTime, group.UpdatedAt)
	})

	t.Run("should do nothing and return no error when target user is not into the group", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUserID := uuid.New().String()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner}).
			Build()
		originalUpdateTime := group.UpdatedAt

		// when
		err := group.RemoveUser(owner.ID, targetUserID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, originalUpdateTime, group.UpdatedAt)
	})

	t.Run("should return forbidden error when trying to remove owner", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner}).
			Build()

		// when
		err := group.RemoveUser(owner.ID, owner.ID)

		// then
		assert.Error(t, err)
		var forbiddenErr *domain.ForbiddenError
		assert.ErrorAs(t, err, &forbiddenErr)
		assert.EqualError(t, forbiddenErr, "cannot remove group owner")
		assert.Contains(t, group.Users, owner)
	})

	t.Run("should return forbidden error when requester is not owner or target user", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		requester := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, targetUser}).
			Build()

		// when
		err := group.RemoveUser(requester.ID, targetUser.ID)

		// then
		assert.Error(t, err)
		var forbiddenErr *domain.ForbiddenError
		assert.ErrorAs(t, err, &forbiddenErr)
		assert.EqualError(t, forbiddenErr, "only the group owner can remove other users")
		assert.Contains(t, group.Users, targetUser)
	})
}

func Test_Group_GenerateMatches(t *testing.T) {
	t.Run("should generate matches successfully", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		user3 := build_domain.NewUserBuilder().Build()

		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithUsers([]domain.User{owner, user1, user2, user3}).Build()

		// when
		err := group.GenerateMatches(owner.ID)

		// then
		assert.NoError(t, err)
		assert.Len(t, group.Matches, 4)

		for _, match := range group.Matches {
			assert.NotEqual(t, match.GiverID, match.ReceiverID)
			assert.Contains(t, []string{owner.ID, user1.ID, user2.ID, user3.ID}, match.GiverID)
			assert.Contains(t, []string{owner.ID, user1.ID, user2.ID, user3.ID}, match.ReceiverID)
		}

		// Ensure all users are givers and receivers exactly once
		givers := make(map[string]int)
		receivers := make(map[string]int)
		for _, match := range group.Matches {
			givers[match.GiverID]++
			receivers[match.ReceiverID]++
		}

		assert.Len(t, givers, 4)
		assert.Len(t, receivers, 4)
		for _, count := range givers {
			assert.Equal(t, 1, count)
		}
		for _, count := range receivers {
			assert.Equal(t, 1, count)
		}

		// Ensure the updated time is set
		assert.WithinDuration(t, time.Now(), group.UpdatedAt, time.Second)
	})

	t.Run("should return an error when requester is not group owner", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		user3 := build_domain.NewUserBuilder().Build()

		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithUsers([]domain.User{owner, user1, user2, user3}).Build()
		requesterID := "some-other-user-id"

		// when
		err := group.GenerateMatches(requesterID)

		// then
		assert.Error(t, err)
		var forbiddenError *domain.ForbiddenError
		assert.ErrorAs(t, err, &forbiddenError)
		assert.EqualError(t, forbiddenError, "only the group owner can generate matches")
		assert.Empty(t, group.Matches)
	})

	t.Run("should return an error when group has less than 3 users", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()

		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithUsers([]domain.User{owner, user1}).Build()

		// when
		err := group.GenerateMatches(owner.ID)

		// then
		assert.Error(t, err)
		var conflictError *domain.ConflictError
		assert.ErrorAs(t, err, &conflictError)
		assert.EqualError(t, conflictError, "group must have at least 3 users to generate matches")
		assert.Empty(t, group.Matches)
	})
}

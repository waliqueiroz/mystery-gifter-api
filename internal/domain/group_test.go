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
		assert.Equal(t, domain.GroupStatusOpen, group.Status)
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

	t.Run("should return conflict error when group is not open", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithStatus(domain.GroupStatusMatched).Build()

		// when
		err := group.AddUser(owner.ID, targetUser)

		// then
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
		assert.EqualError(t, conflictErr, "group is not open for registration, contact the group owner to reopen the group")
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

	t.Run("should return conflict error when group is not open for removal", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		targetUser := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, targetUser}).
			WithStatus(domain.GroupStatusMatched).Build()

		// when
		err := group.RemoveUser(owner.ID, targetUser.ID)

		// then
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
		assert.EqualError(t, conflictErr, "group is not open for removal, contact the group owner to reopen the group")
		assert.Contains(t, group.Users, targetUser)
	})
}

func Test_Group_Reopen(t *testing.T) {
	t.Run("should reopen a matched group successfully", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusMatched).
			WithMatches([]domain.Match{{GiverID: owner.ID, ReceiverID: user1.ID}}).Build()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Reopen(owner.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, domain.GroupStatusOpen, group.Status)
		assert.Empty(t, group.Matches)
		assert.NotEqual(t, originalUpdatedAt, group.UpdatedAt)
	})

	t.Run("should return forbidden error when requester is not owner", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusMatched).
			WithMatches([]domain.Match{{GiverID: owner.ID, ReceiverID: user1.ID}}).Build()
		requesterID := uuid.New().String()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Reopen(requesterID)

		// then
		assert.Error(t, err)
		var forbiddenErr *domain.ForbiddenError
		assert.ErrorAs(t, err, &forbiddenErr)
		assert.EqualError(t, forbiddenErr, "only the group owner can reopen the group")
		assert.Equal(t, domain.GroupStatusMatched, group.Status)
		assert.NotEmpty(t, group.Matches)
		assert.Equal(t, originalUpdatedAt, group.UpdatedAt)
	})

	t.Run("should return conflict error when group is already open", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusOpen).Build()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Reopen(owner.ID)

		// then
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
		assert.EqualError(t, conflictErr, "group is already open")
		assert.Equal(t, domain.GroupStatusOpen, group.Status)
		assert.Equal(t, originalUpdatedAt, group.UpdatedAt)
	})

	t.Run("should return conflict error when group is archived", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusArchived).Build()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Reopen(owner.ID)

		// then
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
		assert.EqualError(t, conflictErr, "group is archived and cannot be reopened")
		assert.Equal(t, domain.GroupStatusArchived, group.Status)
		assert.Equal(t, originalUpdatedAt, group.UpdatedAt)
	})
}

func Test_Group_Archive(t *testing.T) {
	t.Run("should archive an open group successfully", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusOpen).Build()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Archive(owner.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, domain.GroupStatusArchived, group.Status)
		assert.NotEqual(t, originalUpdatedAt, group.UpdatedAt)
	})

	t.Run("should archive a matched group successfully", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusMatched).Build()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Archive(owner.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, domain.GroupStatusArchived, group.Status)
		assert.NotEqual(t, originalUpdatedAt, group.UpdatedAt)
	})

	t.Run("should return forbidden error when requester is not owner", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusOpen).Build()
		requesterID := uuid.New().String()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Archive(requesterID)

		// then
		assert.Error(t, err)
		var forbiddenErr *domain.ForbiddenError
		assert.ErrorAs(t, err, &forbiddenErr)
		assert.EqualError(t, forbiddenErr, "only the group owner can archive the group")
		assert.Equal(t, domain.GroupStatusOpen, group.Status)
		assert.Equal(t, originalUpdatedAt, group.UpdatedAt)
	})

	t.Run("should return conflict error when group is already archived", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().
			WithOwnerID(owner.ID).
			WithUsers([]domain.User{owner, user1}).
			WithStatus(domain.GroupStatusArchived).Build()
		originalUpdatedAt := group.UpdatedAt

		// when
		err := group.Archive(owner.ID)

		// then
		assert.Error(t, err)
		var conflictErr *domain.ConflictError
		assert.ErrorAs(t, err, &conflictErr)
		assert.EqualError(t, conflictErr, "group is already archived")
		assert.Equal(t, domain.GroupStatusArchived, group.Status)
		assert.Equal(t, originalUpdatedAt, group.UpdatedAt)
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

	t.Run("should return an error when group is not open for matches", func(t *testing.T) {
		// given
		owner := build_domain.NewUserBuilder().Build()
		user1 := build_domain.NewUserBuilder().Build()

		group := build_domain.NewGroupBuilder().WithOwnerID(owner.ID).WithUsers([]domain.User{owner, user1}).WithStatus(domain.GroupStatusMatched).Build()

		// when
		err := group.GenerateMatches(owner.ID)

		// then
		assert.Error(t, err)
		var conflictError *domain.ConflictError
		assert.ErrorAs(t, err, &conflictError)
		assert.EqualError(t, conflictError, "group is not open for matches")
		assert.Empty(t, group.Matches)
	})
}

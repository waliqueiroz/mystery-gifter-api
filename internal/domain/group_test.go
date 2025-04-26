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
		messages := validationErr.Details().(validator.ValidationErrors)
		assert.Len(t, messages, 1)
		assert.Equal(t, "Name", messages[0].Field)
		assert.Equal(t, "Name is a required field", messages[0].Error)
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

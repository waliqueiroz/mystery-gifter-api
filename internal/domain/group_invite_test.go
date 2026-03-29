package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
	"go.uber.org/mock/gomock"
)

func Test_NewGroupInvite(t *testing.T) {
	t.Run("should create a new group invite successfully", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		expiration := 24 * time.Hour
		generatedID := uuid.New().String()
		now := time.Now()

		mockCtrl := gomock.NewController(t)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(generatedID, nil)

		// when
		groupInvite, err := domain.NewGroupInvite(mockedIdentityGenerator, groupID, expiration)

		// then
		assert.NoError(t, err)
		assert.Equal(t, generatedID, groupInvite.ID)
		assert.Equal(t, groupID, groupInvite.GroupID)
		assert.WithinDuration(t, now.Add(expiration), groupInvite.ExpiresAt, time.Second)
		assert.WithinDuration(t, now, groupInvite.CreatedAt, time.Second)
	})

	t.Run("should return error when identity generator fails", func(t *testing.T) {
		// given
		groupID := uuid.New().String()
		expiration := 24 * time.Hour

		mockCtrl := gomock.NewController(t)
		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return("", assert.AnError)

		// when
		groupInvite, err := domain.NewGroupInvite(mockedIdentityGenerator, groupID, expiration)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, groupInvite)
	})
}

func Test_GroupInvite_IsExpired(t *testing.T) {
	t.Run("should return false when invite is not expired", func(t *testing.T) {
		// given
		groupInvite := domain.GroupInvite{
			ID:        uuid.New().String(),
			GroupID:   uuid.New().String(),
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		// when
		result := groupInvite.IsExpired()

		// then
		assert.False(t, result)
	})

	t.Run("should return true when invite is expired", func(t *testing.T) {
		// given
		groupInvite := domain.GroupInvite{
			ID:        uuid.New().String(),
			GroupID:   uuid.New().String(),
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}

		// when
		result := groupInvite.IsExpired()

		// then
		assert.True(t, result)
	})
}

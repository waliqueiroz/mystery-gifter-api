package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres/build_postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres/mock_postgres"
	"go.uber.org/mock/gomock"
)

func Test_groupInviteRepository_Create(t *testing.T) {
	t.Run("should create group invite successfully", func(t *testing.T) {
		// given
		groupInvite := build_domain.NewGroupInviteBuilder().Build()
		insertQuery := "INSERT INTO group_invites (id,group_id,expires_at,created_at) VALUES ($1,$2,$3,$4)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().ExecContext(gomock.Any(), insertQuery, groupInvite.ID, groupInvite.GroupID, groupInvite.ExpiresAt, groupInvite.CreatedAt).Return(nil, nil)

		groupInviteRepository := postgres.NewGroupInviteRepository(mockedDB)

		// when
		err := groupInviteRepository.Create(context.Background(), groupInvite)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when exec fails", func(t *testing.T) {
		// given
		groupInvite := build_domain.NewGroupInviteBuilder().Build()
		insertQuery := "INSERT INTO group_invites (id,group_id,expires_at,created_at) VALUES ($1,$2,$3,$4)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().ExecContext(gomock.Any(), insertQuery, groupInvite.ID, groupInvite.GroupID, groupInvite.ExpiresAt, groupInvite.CreatedAt).Return(nil, assert.AnError)

		groupInviteRepository := postgres.NewGroupInviteRepository(mockedDB)

		// when
		err := groupInviteRepository.Create(context.Background(), groupInvite)

		// then
		assert.Error(t, err)
	})
}

func Test_groupInviteRepository_GetByID(t *testing.T) {
	t.Run("should get group invite by id successfully", func(t *testing.T) {
		// given
		pgGroupInvite := build_postgres.NewGroupInviteBuilder().Build()
		selectQuery := "SELECT * FROM group_invites WHERE id = $1"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectQuery, pgGroupInvite.ID).
			DoAndReturn(func(ctx context.Context, dest any, query string, args ...any) error {
				ptr := dest.(*postgres.GroupInvite)
				*ptr = pgGroupInvite
				return nil
			})

		groupInviteRepository := postgres.NewGroupInviteRepository(mockedDB)

		// when
		result, err := groupInviteRepository.GetByID(context.Background(), pgGroupInvite.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, pgGroupInvite.ID, result.ID)
		assert.Equal(t, pgGroupInvite.GroupID, result.GroupID)
		assert.Equal(t, pgGroupInvite.ExpiresAt, result.ExpiresAt)
		assert.Equal(t, pgGroupInvite.CreatedAt, result.CreatedAt)
	})

	t.Run("should return not found error when invite does not exist", func(t *testing.T) {
		// given
		inviteID := "some-id"
		selectQuery := "SELECT * FROM group_invites WHERE id = $1"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectQuery, inviteID).Return(sql.ErrNoRows)

		groupInviteRepository := postgres.NewGroupInviteRepository(mockedDB)

		// when
		result, err := groupInviteRepository.GetByID(context.Background(), inviteID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var notFoundErr *domain.ResourceNotFoundError
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return error when get fails", func(t *testing.T) {
		// given
		inviteID := "some-id"
		selectQuery := "SELECT * FROM group_invites WHERE id = $1"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectQuery, inviteID).Return(assert.AnError)

		groupInviteRepository := postgres.NewGroupInviteRepository(mockedDB)

		// when
		result, err := groupInviteRepository.GetByID(context.Background(), inviteID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

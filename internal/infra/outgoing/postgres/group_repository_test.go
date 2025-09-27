package postgres_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres/build_postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres/mock_postgres"
	"go.uber.org/mock/gomock"
)

func Test_groupRepository_Create(t *testing.T) {
	t.Run("should create group with one user successfully", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		groupInsertQuery := "INSERT INTO groups (id,name,status,owner_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		groupUsersInsertQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)
		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupInsertQuery, group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupUsersInsertQuery, group.ID, group.Users[0].ID, group.CreatedAt).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(nil)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.NoError(t, err)
	})

	t.Run("should create group with more than one user successfully", func(t *testing.T) {
		// given
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithUsers([]domain.User{user1, user2}).Build()

		groupInsertQuery := "INSERT INTO groups (id,name,status,owner_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		groupUsersInsertQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3),($4,$5,$6)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)
		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupInsertQuery, group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(
			gomock.Any(),
			groupUsersInsertQuery,
			group.ID, group.Users[0].ID, group.CreatedAt,
			group.ID, group.Users[1].ID, group.CreatedAt,
		).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(nil)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.NoError(t, err)
	})

	t.Run("should create group with matches successfully", func(t *testing.T) {
		// given
		match1 := build_domain.NewMatchBuilder().Build()
		match2 := build_domain.NewMatchBuilder().Build()
		group := build_domain.NewGroupBuilder().WithMatches([]domain.Match{match1, match2}).Build()

		groupInsertQuery := "INSERT INTO groups (id,name,status,owner_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		groupUsersInsertQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		groupMatchesInsertQuery := "INSERT INTO group_matches (group_id,giver_id,receiver_id,created_at) VALUES ($1,$2,$3,$4),($5,$6,$7,$8)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)
		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupInsertQuery, group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupUsersInsertQuery, group.ID, group.Users[0].ID, group.CreatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(
			gomock.Any(),
			groupMatchesInsertQuery,
			group.ID, group.Matches[0].GiverID, group.Matches[0].ReceiverID, group.CreatedAt,
			group.ID, group.Matches[1].GiverID, group.Matches[1].ReceiverID, group.CreatedAt,
		).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(nil)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when fail to insert group matches", func(t *testing.T) {
		// given
		match1 := build_domain.NewMatchBuilder().Build()
		group := build_domain.NewGroupBuilder().WithMatches([]domain.Match{match1}).Build()

		groupInsertQuery := "INSERT INTO groups (id,name,status,owner_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		groupUsersInsertQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		groupMatchesInsertQuery := "INSERT INTO group_matches (group_id,giver_id,receiver_id,created_at) VALUES ($1,$2,$3,$4)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupInsertQuery, group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupUsersInsertQuery, group.ID, group.Users[0].ID, group.CreatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(
			gomock.Any(),
			groupMatchesInsertQuery,
			group.ID, group.Matches[0].GiverID, group.Matches[0].ReceiverID, group.CreatedAt,
		).Return(nil, assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error inserting group matches")
	})

	t.Run("should return error when fail to begin transaction", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(nil, assert.AnError)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error beginning transaction")
	})

	t.Run("should return conflict error when group name already exists", func(t *testing.T) {
		// given
		postgresUniqueViolationError := &pq.Error{Code: pq.ErrorCode("23505")}

		group := build_domain.NewGroupBuilder().Build()
		groupInsertQuery := "INSERT INTO groups (id,name,status,owner_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupInsertQuery, group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).Return(nil, postgresUniqueViolationError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.Error(t, err)
		var expectedError *domain.ConflictError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, err, "you already have a group with this name")
	})

	t.Run("should return error when fail to insert group users", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		groupInsertQuery := "INSERT INTO groups (id,name,status,owner_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		groupUsersInsertQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupInsertQuery, group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupUsersInsertQuery, group.ID, group.Users[0].ID, group.CreatedAt).Return(nil, assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error inserting group users")
	})

	t.Run("should return error when fail to commit transaction", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		groupInsertQuery := "INSERT INTO groups (id,name,status,owner_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		groupUsersInsertQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupInsertQuery, group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), groupUsersInsertQuery, group.ID, group.Users[0].ID, group.CreatedAt).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Create(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error committing transaction")
	})
}

func Test_groupRepository_Update(t *testing.T) {
	t.Run("should update group with one user successfully", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		insertUsersQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		deleteMatchesQuery := "DELETE FROM group_matches WHERE group_id = $1"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), insertUsersQuery, group.ID, group.Users[0].ID, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteMatchesQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(nil)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.NoError(t, err)
	})

	t.Run("should update group with more than one user successfully", func(t *testing.T) {
		// given
		user1 := build_domain.NewUserBuilder().Build()
		user2 := build_domain.NewUserBuilder().Build()
		group := build_domain.NewGroupBuilder().WithUsers([]domain.User{user1, user2}).Build()

		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		insertUsersQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3),($4,$5,$6)"
		deleteMatchesQuery := "DELETE FROM group_matches WHERE group_id = $1"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(
			gomock.Any(),
			insertUsersQuery,
			group.ID, group.Users[0].ID, group.UpdatedAt,
			group.ID, group.Users[1].ID, group.UpdatedAt,
		).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteMatchesQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(nil)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.NoError(t, err)
	})

	t.Run("should update group with matches successfully", func(t *testing.T) {
		// given
		match1 := build_domain.NewMatchBuilder().Build()
		match2 := build_domain.NewMatchBuilder().Build()
		group := build_domain.NewGroupBuilder().WithMatches([]domain.Match{match1, match2}).Build()

		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		insertUsersQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		deleteMatchesQuery := "DELETE FROM group_matches WHERE group_id = $1"
		insertMatchesQuery := "INSERT INTO group_matches (group_id,giver_id,receiver_id,created_at) VALUES ($1,$2,$3,$4),($5,$6,$7,$8)"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), insertUsersQuery, group.ID, group.Users[0].ID, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteMatchesQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(
			gomock.Any(),
			insertMatchesQuery,
			group.ID, group.Matches[0].GiverID, group.Matches[0].ReceiverID, group.UpdatedAt,
			group.ID, group.Matches[1].GiverID, group.Matches[1].ReceiverID, group.UpdatedAt,
		).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(nil)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when fail to begin transaction", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(nil, assert.AnError)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error beginning transaction")
	})

	t.Run("should return not found error when group does not exist", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		result := driver.RowsAffected(0)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		var expectedError *domain.ResourceNotFoundError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, err, "group not found")
	})

	t.Run("should return conflict error when group name already exists", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		postgresUniqueViolationError := &pq.Error{Code: pq.ErrorCode("23505")}

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(nil, postgresUniqueViolationError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		var expectedError *domain.ConflictError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, err, "you already have a group with this name")
	})

	t.Run("should return error when fail to update group", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(nil, assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error updating group")
	})

	t.Run("should return error when fail to delete group users", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error deleting group users")
	})

	t.Run("should return error when fail to insert group users", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		insertUsersQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), insertUsersQuery, group.ID, group.Users[0].ID, group.UpdatedAt).Return(nil, assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error inserting group users")
	})

	t.Run("should return error when fail to commit transaction", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		insertUsersQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		deleteMatchesQuery := "DELETE FROM group_matches WHERE group_id = $1"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), insertUsersQuery, group.ID, group.Users[0].ID, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteMatchesQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().Commit().Return(assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error committing transaction")
	})

	t.Run("should return error when fail to delete group matches", func(t *testing.T) {
		// given
		group := build_domain.NewGroupBuilder().Build()
		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		insertUsersQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		deleteMatchesQuery := "DELETE FROM group_matches WHERE group_id = $1"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), insertUsersQuery, group.ID, group.Users[0].ID, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteMatchesQuery, group.ID).Return(nil, assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error deleting group matches")
	})

	t.Run("should return error when fail to insert new group matches", func(t *testing.T) {
		// given
		match1 := build_domain.NewMatchBuilder().Build()
		group := build_domain.NewGroupBuilder().WithMatches([]domain.Match{match1}).Build()

		updateGroupQuery := "UPDATE groups SET name = $1, status = $2, updated_at = $3 WHERE id = $4"
		deleteUsersQuery := "DELETE FROM group_users WHERE group_id = $1"
		insertUsersQuery := "INSERT INTO group_users (group_id,user_id,created_at) VALUES ($1,$2,$3)"
		deleteMatchesQuery := "DELETE FROM group_matches WHERE group_id = $1"
		insertMatchesQuery := "INSERT INTO group_matches (group_id,giver_id,receiver_id,created_at) VALUES ($1,$2,$3,$4)"
		result := driver.RowsAffected(1)

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedTx := mock_postgres.NewMockTX(mockCtrl)

		mockedDB.EXPECT().BeginTxx(gomock.Any(), nil).Return(mockedTx, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), updateGroupQuery, group.Name, group.Status, group.UpdatedAt, group.ID).Return(result, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteUsersQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), insertUsersQuery, group.ID, group.Users[0].ID, group.UpdatedAt).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(gomock.Any(), deleteMatchesQuery, group.ID).Return(nil, nil)
		mockedTx.EXPECT().ExecContext(
			gomock.Any(),
			insertMatchesQuery,
			group.ID, group.Matches[0].GiverID, group.Matches[0].ReceiverID, group.UpdatedAt,
		).Return(nil, assert.AnError)
		mockedTx.EXPECT().Rollback().Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		err := groupRepository.Update(context.Background(), group)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error inserting group matches")
	})
}

func Test_groupRepository_GetByID(t *testing.T) {
	t.Run("should get group by id successfully", func(t *testing.T) {
		// given
		expectedUser1 := build_domain.NewUserBuilder().Build()
		expectedUser2 := build_domain.NewUserBuilder().Build()
		expectedGroup := build_domain.NewGroupBuilder().WithUsers([]domain.User{expectedUser1, expectedUser2}).WithMatches([]domain.Match{}).Build()
		selectGroupQuery := "SELECT g.* FROM groups g WHERE g.id = $1"
		selectUsersQuery := "SELECT u.* FROM users u JOIN group_users gu ON gu.user_id = u.id WHERE gu.group_id = $1"
		selectMatchesQuery := "SELECT giver_id, receiver_id FROM group_matches WHERE group_id = $1"

		group := build_postgres.NewGroupBuilder().
			WithID(expectedGroup.ID).
			WithName(expectedGroup.Name).
			WithStatus(string(expectedGroup.Status)).
			WithOwnerID(expectedGroup.OwnerID).
			WithCreatedAt(expectedGroup.CreatedAt).
			WithUpdatedAt(expectedGroup.UpdatedAt).
			Build()

		user1 := build_postgres.NewUserBuilder().
			WithID(expectedUser1.ID).
			WithName(expectedUser1.Name).
			WithSurname(expectedUser1.Surname).
			WithEmail(expectedUser1.Email).
			WithPassword(expectedUser1.Password).
			WithCreatedAt(expectedUser1.CreatedAt).
			WithUpdatedAt(expectedUser1.UpdatedAt).
			Build()

		user2 := build_postgres.NewUserBuilder().
			WithID(expectedUser2.ID).
			WithName(expectedUser2.Name).
			WithSurname(expectedUser2.Surname).
			WithEmail(expectedUser2.Email).
			WithPassword(expectedUser2.Password).
			WithCreatedAt(expectedUser2.CreatedAt).
			WithUpdatedAt(expectedUser2.UpdatedAt).
			Build()

		users := []postgres.User{user1, user2}
		var matches []postgres.Match // Initialize an empty slice of matches

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectGroupQuery, expectedGroup.ID).SetArg(1, group).Return(nil)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), selectUsersQuery, expectedGroup.ID).SetArg(1, users).Return(nil)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), selectMatchesQuery, expectedGroup.ID).SetArg(1, matches).Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		result, err := groupRepository.GetByID(context.Background(), expectedGroup.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedGroup, result)
	})

	t.Run("should get group by id with matches successfully", func(t *testing.T) {
		// given
		expectedUser1 := build_domain.NewUserBuilder().Build()
		expectedMatch1 := build_domain.NewMatchBuilder().Build()
		expectedMatch2 := build_domain.NewMatchBuilder().Build()
		expectedGroup := build_domain.NewGroupBuilder().WithUsers([]domain.User{expectedUser1}).WithMatches([]domain.Match{expectedMatch1, expectedMatch2}).Build()

		selectGroupQuery := "SELECT g.* FROM groups g WHERE g.id = $1"
		selectUsersQuery := "SELECT u.* FROM users u JOIN group_users gu ON gu.user_id = u.id WHERE gu.group_id = $1"
		selectMatchesQuery := "SELECT giver_id, receiver_id FROM group_matches WHERE group_id = $1"

		group := build_postgres.NewGroupBuilder().
			WithID(expectedGroup.ID).
			WithName(expectedGroup.Name).
			WithOwnerID(expectedGroup.OwnerID).
			WithCreatedAt(expectedGroup.CreatedAt).
			WithUpdatedAt(expectedGroup.UpdatedAt).
			Build()

		user1 := build_postgres.NewUserBuilder().
			WithID(expectedUser1.ID).
			WithName(expectedUser1.Name).
			WithSurname(expectedUser1.Surname).
			WithEmail(expectedUser1.Email).
			WithPassword(expectedUser1.Password).
			WithCreatedAt(expectedUser1.CreatedAt).
			WithUpdatedAt(expectedUser1.UpdatedAt).
			Build()

		users := []postgres.User{user1}

		match1 := build_postgres.NewMatchBuilder().
			WithGiverID(expectedMatch1.GiverID).
			WithReceiverID(expectedMatch1.ReceiverID).
			Build()
		match2 := build_postgres.NewMatchBuilder().
			WithGiverID(expectedMatch2.GiverID).
			WithReceiverID(expectedMatch2.ReceiverID).
			Build()

		matches := []postgres.Match{match1, match2}

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectGroupQuery, expectedGroup.ID).SetArg(1, group).Return(nil)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), selectUsersQuery, expectedGroup.ID).SetArg(1, users).Return(nil)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), selectMatchesQuery, expectedGroup.ID).SetArg(1, matches).Return(nil)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		result, err := groupRepository.GetByID(context.Background(), expectedGroup.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedGroup, result)
	})

	t.Run("should return not found error when group does not exist", func(t *testing.T) {
		// given
		groupID := "550e8400-e29b-41d4-a716-446655440000"
		selectGroupQuery := "SELECT g.* FROM groups g WHERE g.id = $1"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectGroupQuery, groupID).Return(sql.ErrNoRows)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		result, err := groupRepository.GetByID(context.Background(), groupID)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		var expectedError *domain.ResourceNotFoundError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, err, "group not found")
	})

	t.Run("should return error when fail to get group users", func(t *testing.T) {
		// given
		expectedGroup := build_domain.NewGroupBuilder().Build()
		selectGroupQuery := "SELECT g.* FROM groups g WHERE g.id = $1"
		selectUsersQuery := "SELECT u.* FROM users u JOIN group_users gu ON gu.user_id = u.id WHERE gu.group_id = $1"

		group := build_postgres.NewGroupBuilder().
			WithID(expectedGroup.ID).
			WithName(expectedGroup.Name).
			WithOwnerID(expectedGroup.OwnerID).
			WithCreatedAt(expectedGroup.CreatedAt).
			WithUpdatedAt(expectedGroup.UpdatedAt).
			Build()

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)

		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectGroupQuery, expectedGroup.ID).SetArg(1, group).Return(nil)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), selectUsersQuery, expectedGroup.ID).Return(assert.AnError)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		result, err := groupRepository.GetByID(context.Background(), expectedGroup.ID)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "error getting group users")
	})

	t.Run("should return error when fail to get group", func(t *testing.T) {
		// given
		groupID := "550e8400-e29b-41d4-a716-446655440000"
		selectGroupQuery := "SELECT g.* FROM groups g WHERE g.id = $1"

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectGroupQuery, groupID).Return(assert.AnError)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		result, err := groupRepository.GetByID(context.Background(), groupID)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "error getting group")
	})

	t.Run("should return error when fail to get group matches", func(t *testing.T) {
		// given
		expectedGroup := build_domain.NewGroupBuilder().Build()
		selectGroupQuery := "SELECT g.* FROM groups g WHERE g.id = $1"
		selectUsersQuery := "SELECT u.* FROM users u JOIN group_users gu ON gu.user_id = u.id WHERE gu.group_id = $1"
		selectMatchesQuery := "SELECT giver_id, receiver_id FROM group_matches WHERE group_id = $1"

		group := build_postgres.NewGroupBuilder().
			WithID(expectedGroup.ID).
			WithName(expectedGroup.Name).
			WithOwnerID(expectedGroup.OwnerID).
			WithCreatedAt(expectedGroup.CreatedAt).
			WithUpdatedAt(expectedGroup.UpdatedAt).
			Build()

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), selectGroupQuery, expectedGroup.ID).SetArg(1, group).Return(nil)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), selectUsersQuery, expectedGroup.ID).Return(nil)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), selectMatchesQuery, expectedGroup.ID).Return(assert.AnError)

		groupRepository := postgres.NewGroupRepository(mockedDB)

		// when
		result, err := groupRepository.GetByID(context.Background(), expectedGroup.ID)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "error getting group matches")
	})
}

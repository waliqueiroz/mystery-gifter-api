package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/postgres/build_postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/postgres/mock_postgres"
	"go.uber.org/mock/gomock"
)

func Test_userRepository_Create(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()
		query := `INSERT INTO users (id,name,surname,email,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().ExecContext(gomock.Any(), query, user.ID, user.Name, user.Surname, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Return(nil, nil)

		userRespository := postgres.NewUserRepository(mockedDB)

		// when
		err := userRespository.Create(context.Background(), user)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return a conflict error when email is already registered", func(t *testing.T) {
		// given
		postgresUniqueViolationError := &pq.Error{Code: pq.ErrorCode("23505")}

		user := build_domain.NewUserBuilder().Build()
		query := `INSERT INTO users (id,name,surname,email,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().ExecContext(gomock.Any(), query, user.ID, user.Name, user.Surname, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Return(nil, postgresUniqueViolationError)

		userRespository := postgres.NewUserRepository(mockedDB)

		// when
		err := userRespository.Create(context.Background(), user)

		// then
		assert.Error(t, err)
		var expectedError *domain.ConflictError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "the email is already registered")
	})

	t.Run("should fail when db return any other error", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()
		query := `INSERT INTO users (id,name,surname,email,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().ExecContext(gomock.Any(), query, user.ID, user.Name, user.Surname, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Return(nil, assert.AnError)

		userRespository := postgres.NewUserRepository(mockedDB)

		// when
		err := userRespository.Create(context.Background(), user)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_userRepository_GetByID(t *testing.T) {
	t.Run("should get user by id successfully", func(t *testing.T) {
		// given
		userID := uuid.New().String()
		now := time.Now().UTC()
		dbUser := build_postgres.NewUserBuilder().WithID(userID).WithCreatedAt(now).WithUpdatedAt(now).Build()
		domainUser := build_domain.NewUserBuilder().WithID(userID).WithCreatedAt(now).WithUpdatedAt(now).Build()
		query := `SELECT * FROM users WHERE id = $1`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), query, userID).SetArg(1, dbUser).Return(nil)

		userRespository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRespository.GetByID(context.Background(), userID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, domainUser, *result)
	})

	t.Run("should return a resource not found error when user is not found", func(t *testing.T) {
		// given
		userID := uuid.New().String()
		query := `SELECT * FROM users WHERE id = $1`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), query, userID).Return(sql.ErrNoRows)

		userRespository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRespository.GetByID(context.Background(), userID)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ResourceNotFoundError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "user not found")
	})

	t.Run("should fail when db return any other error", func(t *testing.T) {
		// given
		userID := uuid.New().String()
		query := `SELECT * FROM users WHERE id = $1`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), query, userID).Return(assert.AnError)

		userRespository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRespository.GetByID(context.Background(), userID)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})
}

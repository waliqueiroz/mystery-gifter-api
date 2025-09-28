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
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres/build_postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres/mock_postgres"
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

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		err := userRepository.Create(context.Background(), user)

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

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		err := userRepository.Create(context.Background(), user)

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

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		err := userRepository.Create(context.Background(), user)

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

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.GetByID(context.Background(), userID)

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

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.GetByID(context.Background(), userID)

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

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.GetByID(context.Background(), userID)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})
}

func Test_userRepository_GetByEmail(t *testing.T) {
	t.Run("should get user by email successfully", func(t *testing.T) {
		// given
		userID := uuid.New().String()
		email := "test@mail.com"
		now := time.Now().UTC()
		dbUser := build_postgres.NewUserBuilder().WithID(userID).WithEmail(email).WithCreatedAt(now).WithUpdatedAt(now).Build()
		domainUser := build_domain.NewUserBuilder().WithID(userID).WithEmail(email).WithCreatedAt(now).WithUpdatedAt(now).Build()
		query := `SELECT * FROM users WHERE email = $1`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), query, email).SetArg(1, dbUser).Return(nil)

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.GetByEmail(context.Background(), email)

		// then
		assert.NoError(t, err)
		assert.Equal(t, domainUser, *result)
	})

	t.Run("should return a resource not found error when user is not found", func(t *testing.T) {
		// given
		email := "test@mail.com"
		query := `SELECT * FROM users WHERE email = $1`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), query, email).Return(sql.ErrNoRows)

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.GetByEmail(context.Background(), email)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ResourceNotFoundError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "user not found")
	})

	t.Run("should fail when db return any other error", func(t *testing.T) {
		// given
		email := "test@mail.com"
		query := `SELECT * FROM users WHERE email = $1`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), query, email).Return(assert.AnError)

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.GetByEmail(context.Background(), email)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})
}

func Test_userRepository_Search(t *testing.T) {
	t.Run("should search users successfully with filters", func(t *testing.T) {
		// given
		name := "John"
		surname := "Doe"
		email := "john@mail.com"
		limit := 10
		offset := 0
		sortBy := "name"
		sortDirection := domain.SortDirectionTypeAsc

		filters := build_domain.NewUserFiltersBuilder().
			WithName(name).
			WithSurname(surname).
			WithEmail(email).
			WithLimit(limit).
			WithOffset(offset).
			WithSortBy(sortBy).
			WithSortDirection(sortDirection).
			Build()

		userID := uuid.New().String()
		now := time.Now().UTC()
		dbUsers := []postgres.User{
			build_postgres.NewUserBuilder().WithID(userID).WithName(name).WithSurname(surname).WithEmail(email).WithCreatedAt(now).WithUpdatedAt(now).Build(),
		}

		domainUsers := []domain.User{
			build_domain.NewUserBuilder().WithID(userID).WithName(name).WithSurname(surname).WithEmail(email).WithCreatedAt(now).WithUpdatedAt(now).Build(),
		}

		expectedSearchResult := &domain.SearchResult[domain.User]{
			Result: domainUsers,
			Paging: domain.Paging{
				Total:  1,
				Limit:  limit,
				Offset: offset,
			},
		}

		searchQuery := `SELECT * FROM users WHERE name ILIKE $1 AND surname ILIKE $2 AND email ILIKE $3 ORDER BY name ASC LIMIT 10 OFFSET 0`
		countQuery := `SELECT COUNT(*) FROM users WHERE name ILIKE $1 AND surname ILIKE $2 AND email ILIKE $3`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), searchQuery, "%"+name+"%", "%"+surname+"%", "%"+email+"%").SetArg(1, dbUsers).Return(nil)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), countQuery, "%"+name+"%", "%"+surname+"%", "%"+email+"%").SetArg(1, 1).Return(nil)

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.Search(context.Background(), filters)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedSearchResult, result)
	})

	t.Run("should search users successfully without filters", func(t *testing.T) {
		// given
		limit := 10
		offset := 0
		sortBy := "created_at"
		sortDirection := domain.SortDirectionTypeDesc

		filters := build_domain.NewUserFiltersBuilder().
			WithLimit(limit).
			WithOffset(offset).
			WithSortBy(sortBy).
			WithSortDirection(sortDirection).
			Build()

		userID := uuid.New().String()
		now := time.Now().UTC()
		dbUsers := []postgres.User{
			build_postgres.NewUserBuilder().WithID(userID).WithCreatedAt(now).WithUpdatedAt(now).Build(),
		}

		domainUsers := []domain.User{
			build_domain.NewUserBuilder().WithID(userID).WithCreatedAt(now).WithUpdatedAt(now).Build(),
		}

		expectedSearchResult := &domain.SearchResult[domain.User]{
			Result: domainUsers,
			Paging: domain.Paging{
				Total:  1,
				Limit:  limit,
				Offset: offset,
			},
		}

		searchQuery := `SELECT * FROM users ORDER BY created_at DESC LIMIT 10 OFFSET 0`
		countQuery := `SELECT COUNT(*) FROM users`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), searchQuery).SetArg(1, dbUsers).Return(nil)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), countQuery).SetArg(1, 1).Return(nil)

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.Search(context.Background(), filters)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedSearchResult, result)
	})

	t.Run("should fail when SelectContext returns error", func(t *testing.T) {
		// given
		limit := 10
		offset := 0
		sortBy := "name"
		sortDirection := domain.SortDirectionTypeAsc

		filters := build_domain.NewUserFiltersBuilder().
			WithLimit(limit).
			WithOffset(offset).
			WithSortBy(sortBy).
			WithSortDirection(sortDirection).
			Build()

		searchQuery := `SELECT * FROM users ORDER BY name ASC LIMIT 10 OFFSET 0`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), searchQuery).Return(assert.AnError)

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.Search(context.Background(), filters)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})

	t.Run("should fail when countUsers returns error", func(t *testing.T) {
		// given
		limit := 10
		offset := 0
		sortBy := "name"
		sortDirection := domain.SortDirectionTypeAsc

		filters := build_domain.NewUserFiltersBuilder().
			WithLimit(limit).
			WithOffset(offset).
			WithSortBy(sortBy).
			WithSortDirection(sortDirection).
			Build()

		dbUsers := []postgres.User{
			build_postgres.NewUserBuilder().Build(),
		}

		searchQuery := `SELECT * FROM users ORDER BY name ASC LIMIT 10 OFFSET 0`
		countQuery := `SELECT COUNT(*) FROM users`

		mockCtrl := gomock.NewController(t)
		mockedDB := mock_postgres.NewMockDB(mockCtrl)
		mockedDB.EXPECT().SelectContext(gomock.Any(), gomock.Any(), searchQuery).SetArg(1, dbUsers).Return(nil)
		mockedDB.EXPECT().GetContext(gomock.Any(), gomock.Any(), countQuery).Return(assert.AnError)

		userRepository := postgres.NewUserRepository(mockedDB)

		// when
		result, err := userRepository.Search(context.Background(), filters)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})
}

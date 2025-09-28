package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
	"go.uber.org/mock/gomock"
)

func Test_userService_Create(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().Create(gomock.Any(), user).Return(nil)

		userService := application.NewUserService(mockedUserRepository)

		// when
		err := userService.Create(context.Background(), user)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return an error when repository fails", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().Create(gomock.Any(), user).Return(assert.AnError)

		userService := application.NewUserService(mockedUserRepository)

		// when
		err := userService.Create(context.Background(), user)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return a validation error when user is invalid", func(t *testing.T) {
		// given
		invalidUser := build_domain.NewUserBuilder().WithName("").Build()

		userService := application.NewUserService(nil)

		// when
		err := userService.Create(context.Background(), invalidUser)

		// then
		assert.Error(t, err)
		var expectedError *domain.ValidationError
		assert.ErrorAs(t, err, &expectedError)
		assert.Equal(t, "validation failed", expectedError.Error())
		errors := expectedError.Details()
		assert.Contains(t, errors, validator.FieldError{Field: "Name", Error: "Name is a required field"})
	})
}

func Test_userService_GetByID(t *testing.T) {
	t.Run("should get user by id successfully", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().GetByID(gomock.Any(), user.ID).Return(&user, nil)

		userService := application.NewUserService(mockedUserRepository)

		// when
		result, err := userService.GetByID(context.Background(), user.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, user, *result)
	})

	t.Run("should return an error when repository fails", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().GetByID(gomock.Any(), user.ID).Return(nil, assert.AnError)

		userService := application.NewUserService(mockedUserRepository)

		// when
		result, err := userService.GetByID(context.Background(), user.ID)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})
}

func Test_userService_Search(t *testing.T) {
	t.Run("should search users successfully", func(t *testing.T) {
		// given
		filters := build_domain.NewUserFiltersBuilder().
			WithName("John").
			WithLimit(10).
			WithOffset(0).
			Build()

		users := []domain.User{
			build_domain.NewUserBuilder().WithName("John").WithSurname("Doe").Build(),
			build_domain.NewUserBuilder().WithName("John").WithSurname("Smith").Build(),
		}

		expectedSearchResult := build_domain.NewSearchResultBuilder[domain.User]().
			WithResult(users).
			WithLimit(10).
			WithOffset(0).
			WithTotal(2).
			Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().Search(gomock.Any(), filters).Return(&expectedSearchResult, nil)

		userService := application.NewUserService(mockedUserRepository)

		// when
		result, err := userService.Search(context.Background(), filters)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedSearchResult, *result)
		assert.Len(t, result.Result, 2)
		assert.Equal(t, 10, result.Paging.Limit)
		assert.Equal(t, 0, result.Paging.Offset)
		assert.Equal(t, 2, result.Paging.Total)
	})

	t.Run("should return an error when repository fails", func(t *testing.T) {
		// given
		filters := build_domain.NewUserFiltersBuilder().Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().Search(gomock.Any(), filters).Return(nil, assert.AnError)

		userService := application.NewUserService(mockedUserRepository)

		// when
		result, err := userService.Search(context.Background(), filters)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})

	t.Run("should return a validation error when filters are invalid", func(t *testing.T) {
		// given
		invalidFilters := build_domain.NewUserFiltersBuilder().
			WithLimit(0). // Invalid: limit must be at least 1
			Build()

		userService := application.NewUserService(nil)

		// when
		result, err := userService.Search(context.Background(), invalidFilters)

		// then
		assert.Error(t, err)
		var expectedError *domain.ValidationError
		assert.ErrorAs(t, err, &expectedError)
		assert.Equal(t, "validation failed", expectedError.Error())
		errors := expectedError.Details()
		assert.Contains(t, errors, validator.FieldError{Field: "Limit", Error: "Limit is a required field"})
		assert.Nil(t, result)
	})

	t.Run("should return empty result when no users match filters", func(t *testing.T) {
		// given
		filters := build_domain.NewUserFiltersBuilder().
			WithName("NonExistentUser").
			Build()

		emptyResult := build_domain.NewSearchResultBuilder[domain.User]().
			WithResult([]domain.User{}).
			WithLimit(10).
			WithOffset(0).
			WithTotal(0).
			Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().Search(gomock.Any(), filters).Return(&emptyResult, nil)

		userService := application.NewUserService(mockedUserRepository)

		// when
		result, err := userService.Search(context.Background(), filters)

		// then
		assert.NoError(t, err)
		assert.Equal(t, emptyResult, *result)
		assert.Empty(t, result.Result)
		assert.Equal(t, 0, result.Paging.Total)
	})
}

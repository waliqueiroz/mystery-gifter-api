package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
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

package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/mock_domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
	"go.uber.org/mock/gomock"
)

func Test_authService_Login(t *testing.T) {
	t.Run("should return a session successfully", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "some_password"
		hashedPassword := "some_hashed_password"
		token := "some_token"
		tokenType := "Bearer"
		sessionDuration := time.Hour
		credentials := build_domain.NewCredentialsBuilder().WithEmail(email).WithPassword(password).Build()
		user := build_domain.NewUserBuilder().WithEmail(email).WithPassword(hashedPassword).Build()

		authSession := build_domain.NewAuthSessionBuilder().WithUser(user).WithAccessToken(token).WithTokenType(tokenType).WithExpiresIn(time.Now().Add(sessionDuration).Unix()).Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().GetByEmail(gomock.Any(), credentials.Email).Return(&user, nil)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Compare(user.Password, credentials.Password).Return(nil)

		mockedTokenManager := mock_domain.NewMockTokenManager(mockCtrl)
		mockedTokenManager.EXPECT().Create(user.ID, gomock.Any()).Return(token, nil)
		mockedTokenManager.EXPECT().GetTokenType().Return(tokenType)

		authService := application.NewAuthService(sessionDuration, mockedUserRepository, mockedPasswordManager, mockedTokenManager)

		// when
		result, err := authService.Login(context.Background(), credentials)

		// then
		assert.NoError(t, err)
		assert.Equal(t, authSession, *result)
	})

	t.Run("should return an error when token creation fails", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "some_password"
		hashedPassword := "some_hashed_password"
		sessionDuration := time.Hour
		credentials := build_domain.NewCredentialsBuilder().WithEmail(email).WithPassword(password).Build()
		user := build_domain.NewUserBuilder().WithEmail(email).WithPassword(hashedPassword).Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().GetByEmail(gomock.Any(), credentials.Email).Return(&user, nil)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Compare(user.Password, credentials.Password).Return(nil)

		mockedTokenManager := mock_domain.NewMockTokenManager(mockCtrl)
		mockedTokenManager.EXPECT().Create(user.ID, gomock.Any()).Return("", assert.AnError)

		authService := application.NewAuthService(sessionDuration, mockedUserRepository, mockedPasswordManager, mockedTokenManager)

		// when
		result, err := authService.Login(context.Background(), credentials)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return an unauthorized error when password comparison fails", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "some_password"
		hashedPassword := "some_hashed_password"
		sessionDuration := time.Hour
		credentials := build_domain.NewCredentialsBuilder().WithEmail(email).WithPassword(password).Build()
		user := build_domain.NewUserBuilder().WithEmail(email).WithPassword(hashedPassword).Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().GetByEmail(gomock.Any(), credentials.Email).Return(&user, nil)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Compare(user.Password, credentials.Password).Return(assert.AnError)

		authService := application.NewAuthService(sessionDuration, mockedUserRepository, mockedPasswordManager, nil)

		// when
		result, err := authService.Login(context.Background(), credentials)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.UnauthorizedError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "invalid credentials")
	})

	t.Run("should return an unauthorized error when it fails to get user by email", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "some_password"
		sessionDuration := time.Hour
		credentials := build_domain.NewCredentialsBuilder().WithEmail(email).WithPassword(password).Build()

		mockCtrl := gomock.NewController(t)
		mockedUserRepository := mock_domain.NewMockUserRepository(mockCtrl)
		mockedUserRepository.EXPECT().GetByEmail(gomock.Any(), credentials.Email).Return(nil, assert.AnError)

		authService := application.NewAuthService(sessionDuration, mockedUserRepository, nil, nil)

		// when
		result, err := authService.Login(context.Background(), credentials)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.UnauthorizedError
		assert.ErrorAs(t, err, &expectedError)
		assert.EqualError(t, expectedError, "invalid credentials")
	})

	t.Run("should return a validation error when credentials are invalid", func(t *testing.T) {
		// given
		sessionDuration := time.Hour
		credentials := build_domain.NewCredentialsBuilder().WithEmail("").WithPassword("").Build()

		authService := application.NewAuthService(sessionDuration, nil, nil, nil)

		// when
		result, err := authService.Login(context.Background(), credentials)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var expectedError *domain.ValidationError
		assert.ErrorAs(t, err, &expectedError)
		assert.Equal(t, "validation failed", expectedError.Error())
		errors := expectedError.Details()
		assert.Len(t, errors, 2)
		assert.Contains(t, errors, validator.FieldError{Field: "Email", Error: "Email is a required field"})
		assert.Contains(t, errors, validator.FieldError{Field: "Password", Error: "Password is a required field"})
	})
}

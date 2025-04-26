package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain/build_domain"
)

func Test_NewCredentials(t *testing.T) {
	t.Run("should create new credentials successfully", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := "12345678"

		// when
		credentials, err := domain.NewCredentials(email, password)

		// then
		assert.NoError(t, err)
		assert.Equal(t, email, credentials.Email)
		assert.Equal(t, password, credentials.Password)
	})

	t.Run("should return validation error when email is empty", func(t *testing.T) {
		// given
		email := ""
		password := "12345678"

		// when
		credentials, err := domain.NewCredentials(email, password)

		// then
		assert.Error(t, err)
		var validationErr *domain.ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Nil(t, credentials)
	})

	t.Run("should return validation error when password is empty", func(t *testing.T) {
		// given
		email := "test@mail.com"
		password := ""

		// when
		credentials, err := domain.NewCredentials(email, password)

		// then
		assert.Error(t, err)
		var validationErr *domain.ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Nil(t, credentials)
	})
}

func Test_NewAuthSession(t *testing.T) {
	t.Run("should create new auth session successfully", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()
		accessToken := "some_token"
		tokenType := "Bearer"
		now := time.Now()
		expiresIn := now.Add(time.Hour).Unix()

		// when
		authSession, err := domain.NewAuthSession(user, accessToken, tokenType, expiresIn)

		// then
		assert.NoError(t, err)
		assert.Equal(t, user, authSession.User)
		assert.Equal(t, accessToken, authSession.AccessToken)
		assert.Equal(t, tokenType, authSession.TokenType)
		assert.Equal(t, expiresIn, authSession.ExpiresIn)
	})

	t.Run("should return validation error when access token is empty", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()
		accessToken := ""
		tokenType := "Bearer"
		expiresIn := time.Now().Add(time.Hour).Unix()

		// when
		authSession, err := domain.NewAuthSession(user, accessToken, tokenType, expiresIn)

		// then
		assert.Error(t, err)
		var validationErr *domain.ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Nil(t, authSession)
	})

	t.Run("should return validation error when token type is empty", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().Build()
		accessToken := "some_token"
		tokenType := ""
		expiresIn := time.Now().Add(time.Hour).Unix()

		// when
		authSession, err := domain.NewAuthSession(user, accessToken, tokenType, expiresIn)

		// then
		assert.Error(t, err)
		var validationErr *domain.ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Nil(t, authSession)
	})

	t.Run("should return validation error when user is invalid", func(t *testing.T) {
		// given
		user := build_domain.NewUserBuilder().WithName("").Build()
		accessToken := "some_token"
		tokenType := "Bearer"
		expiresIn := time.Now().Add(time.Hour).Unix()

		// when
		authSession, err := domain.NewAuthSession(user, accessToken, tokenType, expiresIn)

		// then
		assert.Error(t, err)
		var validationErr *domain.ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Nil(t, authSession)
	})
}

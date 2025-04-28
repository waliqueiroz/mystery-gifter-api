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

func Test_NewUser(t *testing.T) {
	t.Run("should create a new user successfully", func(t *testing.T) {
		// given
		name := "John"
		surname := "Doe"
		email := "john@example.com"
		password := "password123"
		hashedPassword := "hashed_password"
		userID := uuid.New().String()
		now := time.Now()

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(password).Return(hashedPassword, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(userID, nil)

		expectedUser := build_domain.NewUserBuilder().
			WithID(userID).
			WithName(name).
			WithSurname(surname).
			WithEmail(email).
			WithPassword(hashedPassword).
			WithCreatedAt(now).
			WithUpdatedAt(now).
			Build()

		// when
		result, err := domain.NewUser(mockedIdentityGenerator, mockedPasswordManager, name, surname, email, password)
		result.CreatedAt = now
		result.UpdatedAt = now

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, *result)
	})

	t.Run("should return error when password manager fails", func(t *testing.T) {
		// given
		name := "John"
		surname := "Doe"
		email := "john@example.com"
		password := "password123"

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(password).Return("", assert.AnError)

		// when
		result, err := domain.NewUser(nil, mockedPasswordManager, name, surname, email, password)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})

	t.Run("should return error when identity generator fails", func(t *testing.T) {
		// given
		name := "John"
		surname := "Doe"
		email := "john@example.com"
		password := "password123"
		hashedPassword := "hashed_password"

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(password).Return(hashedPassword, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return("", assert.AnError)

		// when
		result, err := domain.NewUser(mockedIdentityGenerator, mockedPasswordManager, name, surname, email, password)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
	})

	t.Run("should return validation error when name is empty", func(t *testing.T) {
		// given
		name := ""
		surname := "Doe"
		email := "john@example.com"
		password := "password123"
		hashedPassword := "hashed_password"
		userID := uuid.New().String()

		mockCtrl := gomock.NewController(t)

		mockedPasswordManager := mock_domain.NewMockPasswordManager(mockCtrl)
		mockedPasswordManager.EXPECT().Hash(password).Return(hashedPassword, nil)

		mockedIdentityGenerator := mock_domain.NewMockIdentityGenerator(mockCtrl)
		mockedIdentityGenerator.EXPECT().Generate().Return(userID, nil)

		// when
		result, err := domain.NewUser(mockedIdentityGenerator, mockedPasswordManager, name, surname, email, password)

		// then
		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *domain.ValidationError
		assert.ErrorAs(t, err, &validationErr)
		errors := validationErr.Details()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors, validator.FieldError{Field: "Name", Error: "Name is a required field"})
	})
}

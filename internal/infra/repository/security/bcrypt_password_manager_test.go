package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/security"
)

func Test_BcryptPasswordManager_Hash(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		// given
		password := "securepassword"
		passwordManager := security.NewBcryptPasswordManager()

		// when
		hashedPassword, err := passwordManager.Hash(password)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("should return an error when password is too long", func(t *testing.T) {
		// given
		password := "012345678901234567890123456789012345678901234567890123456012345678901234567890123456789012345678901234567890123456"
		passwordManager := security.NewBcryptPasswordManager()

		// when
		hashedPassword, err := passwordManager.Hash(password)

		// then
		assert.Empty(t, hashedPassword)
		assert.Error(t, err)
		assert.EqualError(t, err, "bcrypt: password length exceeds 72 bytes")
	})
}

func Test_BcryptPasswordManager_Compare(t *testing.T) {
	t.Run("should compare password successfully", func(t *testing.T) {
		// given
		password := "securepassword"
		passwordManager := security.NewBcryptPasswordManager()
		hashedPassword, err := passwordManager.Hash(password)
		assert.NoError(t, err)

		// when
		err = passwordManager.Compare(hashedPassword, password)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return an error when password does not match", func(t *testing.T) {
		// given
		password := "securepassword"
		otherPassword := "wrongpassword"
		passwordManager := security.NewBcryptPasswordManager()
		hashedPassword, err := passwordManager.Hash(password)
		assert.NoError(t, err)

		// when
		err = passwordManager.Compare(hashedPassword, otherPassword)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "crypto/bcrypt: hashedPassword is not the hash of the given password")
	})

	t.Run("should return an error when hashed password is invalid", func(t *testing.T) {
		// given
		password := "securepassword"
		invalidHashedPassword := "invalidHash"
		passwordManager := security.NewBcryptPasswordManager()

		// when
		err := passwordManager.Compare(invalidHashedPassword, password)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "crypto/bcrypt: hashedSecret too short to be a bcrypted password")
	})
}

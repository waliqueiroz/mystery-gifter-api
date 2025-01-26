package security_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/security"
)

func Test_JWTTokenManager_Create(t *testing.T) {
	t.Run("should create token successfully", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		userID := "some-user-id"
		expiresIn := time.Now().Add(time.Hour).Unix()

		tokenManager := security.NewJWTTokenManager(secretKey)

		// when
		token, err := tokenManager.Create(userID, expiresIn)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, userID, claims["userID"])
		assert.Equal(t, true, claims["authorized"])
		assert.Equal(t, expiresIn, int64(claims["exp"].(float64)))
	})
}

func Test_JWTTokenManager_GetTokenType(t *testing.T) {
	t.Run("should return the correct token type", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		tokenManager := security.NewJWTTokenManager(secretKey)

		// when
		tokenType := tokenManager.GetTokenType()

		// then
		assert.Equal(t, "Bearer", tokenType)
	})
}

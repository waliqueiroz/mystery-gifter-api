package security_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/security"
)

func Test_JWTSessionManager_Create(t *testing.T) {
	t.Run("should create token successfully", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		userID := "some-user-id"
		expiresIn := time.Now().Add(time.Hour).Unix()

		sessionManager := security.NewJWTSessionManager(secretKey)

		// when
		token, err := sessionManager.Create(userID, expiresIn)

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

func Test_JWTSessionManager_GetTokenType(t *testing.T) {
	t.Run("should return the correct token type", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		sessionManager := security.NewJWTSessionManager(secretKey)

		// when
		tokenType := sessionManager.GetTokenType()

		// then
		assert.Equal(t, "Bearer", tokenType)
	})
}

func Test_JWTSessionManager_ExtractUserID(t *testing.T) {
	t.Run("should extract user ID successfully", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		sessionManager := security.NewJWTSessionManager(secretKey)
		expectedUserID := "some-user-id"

		token := jwt.New(jwt.SigningMethodHS256)
		token.Valid = true
		token.Claims = jwt.MapClaims{
			"authorized": true,
			"exp":        time.Now().Add(time.Hour).Unix(),
			"userID":     expectedUserID,
		}

		// when
		userID, err := sessionManager.ExtractUserID(token)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
	})

	t.Run("should return error when token is not a JWT token", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		invalidToken := "invalid-token"

		sessionManager := security.NewJWTSessionManager(secretKey)

		// when
		userID, err := sessionManager.ExtractUserID(invalidToken)

		// then
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.EqualError(t, err, "invalid token")
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		sessionManager := security.NewJWTSessionManager(secretKey)

		invalidToken := jwt.New(jwt.SigningMethodHS256)
		invalidToken.Valid = false

		// when
		userID, err := sessionManager.ExtractUserID(invalidToken)

		// then
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.EqualError(t, err, "invalid token")
	})

	t.Run("should return error when userID claim is missing", func(t *testing.T) {
		// given
		secretKey := "mysecretkey"
		sessionManager := security.NewJWTSessionManager(secretKey)

		token := jwt.New(jwt.SigningMethodHS256)
		token.Valid = true
		token.Claims = jwt.MapClaims{
			"authorized": true,
			"exp":        time.Now().Add(time.Hour).Unix(),
		}

		// when
		userID, err := sessionManager.ExtractUserID(token)

		// then
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.EqualError(t, err, "invalid token")
	})
}

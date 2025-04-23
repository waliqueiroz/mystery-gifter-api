package security

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type JWTTokenManager struct {
	secretKey string
	tokenType string
}

func NewJWTTokenManager(secretKey string) domain.TokenManager {
	return &JWTTokenManager{
		secretKey: secretKey,
		tokenType: "Bearer",
	}
}

func (t *JWTTokenManager) Create(userID string, expiresIn int64) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        expiresIn,
		"userID":     userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}

func (t *JWTTokenManager) GetTokenType() string {
	return t.tokenType
}

func (t *JWTTokenManager) ExtractUserID(token any) (string, error) {
	err := domain.NewUnauthorizedError("invalid token")

	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", err
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", err
	}

	return userID, nil
}

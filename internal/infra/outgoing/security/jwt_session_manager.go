package security

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type JWTSessionManager struct {
	secretKey string
	tokenType string
}

func NewJWTSessionManager(secretKey string) domain.SessionManager {
	return &JWTSessionManager{
		secretKey: secretKey,
		tokenType: "Bearer",
	}
}

func (t *JWTSessionManager) Create(userID string, expiresIn int64) (string, error) {
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

func (t *JWTSessionManager) GetTokenType() string {
	return t.tokenType
}

func (t *JWTSessionManager) ExtractUserID(token any) (string, error) {
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

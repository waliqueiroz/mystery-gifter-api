package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/auth_token_manager.go . AuthTokenManager
//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/password_manager.go . PasswordManager

type PasswordManager interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

type AuthTokenManager interface {
	Create(userID string, expiresIn int64) (string, error)
	GetTokenType() string
	GetAuthUserID(token any) (string, error)
}

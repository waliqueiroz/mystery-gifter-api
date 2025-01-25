package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/password_manager.go . PasswordManager
//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/token_manager.go . TokenManager

type PasswordManager interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

type TokenManager interface {
	Create(userID string, expiresIn int64) (string, error)
	GetTokenType() string
}

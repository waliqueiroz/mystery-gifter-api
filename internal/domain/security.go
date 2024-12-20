package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/password_manager.go . PasswordManager

type PasswordManager interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

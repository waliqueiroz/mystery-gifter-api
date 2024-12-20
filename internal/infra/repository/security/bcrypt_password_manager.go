package security

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordManager struct{}

func NewBcryptPasswordManager() domain.PasswordManager {
	return &BcryptPasswordManager{}
}

// Hash takes a string an d returns a hash of it
func (b *BcryptPasswordManager) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare compares a bcrypt hashed password with its possible plaintext equivalent. Returns nil on success, or an error on failure
func (b *BcryptPasswordManager) Compare(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

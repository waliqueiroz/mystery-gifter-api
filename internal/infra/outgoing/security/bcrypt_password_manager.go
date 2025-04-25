package security

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordManager struct{}

func NewBcryptPasswordManager() domain.PasswordManager {
	return &BcryptPasswordManager{}
}

func (b *BcryptPasswordManager) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptPasswordManager) Compare(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

package security

import "golang.org/x/crypto/bcrypt"

// Hash takes a string and returns a hash of it
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Verify compares a bcrypt hashed password with its possible plaintext equivalent. Returns nil on success, or an error on failure
func Verify(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

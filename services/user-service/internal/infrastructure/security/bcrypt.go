package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements domain.PasswordHasher using bcrypt.
type BcryptHasher struct{}

// NewBcryptHasher creates a new BcryptHasher.
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// Hash hashes a plain-text password using bcrypt.
func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

// Compare compares a hashed password with a plain-text password.
func (h *BcryptHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

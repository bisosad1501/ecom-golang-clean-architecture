package services

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) error
}

type passwordService struct {
	cost int
}

// NewPasswordService creates a new password service
func NewPasswordService() PasswordService {
	return &passwordService{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword hashes a password using bcrypt
func (s *passwordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword checks if a password matches the hashed password
func (s *passwordService) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

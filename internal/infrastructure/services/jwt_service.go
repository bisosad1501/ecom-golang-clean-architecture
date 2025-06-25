package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService implements JWT token generation
type JWTService struct {
	secret string
}

// NewJWTService creates a new JWT service
func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: secret,
	}
}

// GenerateToken generates a JWT token for the given user ID and role
func (s *JWTService) GenerateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

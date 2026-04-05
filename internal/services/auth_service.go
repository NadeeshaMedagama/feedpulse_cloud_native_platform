package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	adminEmail    string
	adminPassword string
	jwtSecret     []byte
}

func NewAuthService(adminEmail, adminPassword, jwtSecret string) *AuthService {
	return &AuthService{adminEmail: adminEmail, adminPassword: adminPassword, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Login(email, password string) (string, error) {
	if email != s.adminEmail || password != s.adminPassword {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":  email,
		"role": "admin",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
}

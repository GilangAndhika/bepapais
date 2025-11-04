package auth

import (
	"context"
	"errors"
	"time"

	"github.com/GilangAndhika/bepapais/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

// NewAuthService membuat instance service baru
func NewAuthService(repo Repository, jwtSecret string) Service {
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Login memvalidasi user dan mengembalikan JWT token
func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	// 1. Cari user di DB
	admin, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if admin == nil {
		return "", ErrUserNotFound
	}

	// 2. Bandingkan password
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		// Password tidak cocok
		return "", ErrInvalidCredentials
	}

	// 3. Buat JWT Token
	token, err := s.generateJWT(admin)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateJWT membuat token baru untuk admin
func (s *service) generateJWT(admin *models.Admin) (string, error) {
	// Buat claims
	claims := jwt.MapClaims{
		"id":       admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 72 jam
		"iat":      time.Now().Unix(),
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret
	t, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
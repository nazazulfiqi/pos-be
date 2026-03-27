// internal/service/auth_service.go
package service

import (
	"errors"
	"os"
	"time"

	"pos-be/internal/dto"
	"pos-be/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignIn(req dto.SignInRequest) (dto.SignInResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) SignIn(req dto.SignInRequest) (dto.SignInResponse, error) {
	// cari user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return dto.SignInResponse{}, errors.New("invalid email or password")
	}

	// cek password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return dto.SignInResponse{}, errors.New("invalid email or password")
	}

	// fetch role slugs for the user and generate JWT
	roleSlugs, err := s.repo.GetRoleSlugsByUser(user.ID)
	if err != nil {
		return dto.SignInResponse{}, err
	}

	// fetch permission slugs for the user
	permissionSlugs, err := s.repo.GetPermissionSlugsByUser(user.ID)
	if err != nil {
		return dto.SignInResponse{}, err
	}

	claims := jwt.MapClaims{
		"user_id":     user.ID,
		"roles":       roleSlugs,
		"permissions": permissionSlugs,
		"tenant_id":   user.TenantID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // expired 24 jam
		"iat":         time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret123" // default dev
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return dto.SignInResponse{}, err
	}

	return dto.SignInResponse{AccessToken: tokenString}, nil
}

package service

import (
	"domainhub/internal/dto"
	"domainhub/internal/models"
	"domainhub/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) error {
	if req.Email == "" || req.Password == "" {
		return errors.New("email and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "USER",
	}
	return s.userRepo.Create(user)
}

package service

import (
	"errors"

	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type AuthService interface {
	Register(username, email, password string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	ValidateToken(token string) (string, string, error)
	GetUser(userID string) (*models.User, error)
}

type authService struct {
	repo       repository.AuthRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo:       repo,
		jwtManager: utils.NewJWTManager(),
	}
}

func (s *authService) Register(username, email, password string) (*models.User, string, error) {
	// Check if user already exists
	_, err := s.repo.GetUserByEmail(email)
	if err == nil {
		return nil, "", ErrUserAlreadyExists
	}

	_, err = s.repo.GetUserByUsername(username)
	if err == nil {
		return nil, "", ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check password
	if !utils.CheckPassword(password, user.Password) {
		return nil, "", ErrInvalidCredentials
	}

	// Generate token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) ValidateToken(token string) (string, string, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return "", "", err
	}

	return claims.UserID, claims.Username, nil
}

func (s *authService) GetUser(userID string) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
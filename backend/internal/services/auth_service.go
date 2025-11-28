package services

import (
	"bobastream/internal/models"
	"bobastream/internal/repositories"
	"bobastream/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register registers a new user
func (s *AuthService) Register(email, username, password string) (*models.User, error) {
	// Check if email exists
	emailExists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("email already registered")
	}

	// Check if username exists
	usernameExists, err := s.userRepo.UsernameExists(username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         models.RoleViewer,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates user and returns tokens
func (s *AuthService) Login(email, password string) (accessToken, refreshToken string, user *models.User, err error) {
	// Find user by email
	user, err = s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, errors.New("invalid credentials")
		}
		return "", "", nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err = utils.GenerateAccessToken(user)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err = utils.GenerateRefreshToken(user)
	if err != nil {
		return "", "", nil, err
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Non-critical error, just log it
		// In production, you'd use proper logging here
	}

	return accessToken, refreshToken, user, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// GetUserByID gets user by ID
func (s *AuthService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// ValidateToken validates access token and returns user
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(claims.UserID)
}
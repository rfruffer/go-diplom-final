package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fylgushev/go-diplom-final/internal/infrastructure/crypto"
	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/fylgushev/go-diplom-final/pkg/proto/auth"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
)

// UserRepository defines the interface for user storage operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByLogin(ctx context.Context, login string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
}

// AuthService implements authentication business logic
type AuthService struct {
	userRepo     UserRepository
	tokenManager *TokenManager
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo UserRepository, tokenManager *TokenManager) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if err := validateCredentials(req.Login, req.Password); err != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByLogin(ctx, req.Login)
	if err == nil && existingUser != nil {
		return &auth.RegisterResponse{
			Success: false,
			Message: "user already exists",
		}, nil
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entities.User{
		ID:           uuid.New().String(),
		Username:     req.Login,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate and return tokens
	tokens, err := s.generateUserTokens(user)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Success:      true,
		Message:      "user registered successfully",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    timestamppb.New(tokens.ExpiresAt),
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := validateCredentials(req.Login, req.Password); err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Get user by login
	user, err := s.userRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: "invalid login or password",
		}, nil
	}

	// Verify password
	if !crypto.VerifyPassword(req.Password, user.PasswordHash) {
		return &auth.LoginResponse{
			Success: false,
			Message: "invalid login or password",
		}, nil
	}

	// Generate tokens
	tokens, err := s.generateUserTokens(user)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Success:      true,
		Message:      "login successful",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    timestamppb.New(tokens.ExpiresAt),
		User:         convertUserToProto(user),
	}, nil
}

// ValidateToken validates a JWT token and returns user information
func (s *AuthService) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &auth.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	claims, err := s.tokenManager.ValidateToken(req.Token)
	if err != nil {
		return &auth.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &auth.ValidateTokenResponse{
		Valid:     true,
		UserId:    claims.UserID,
		Login:     claims.Login,
		ExpiresAt: timestamppb.New(claims.ExpiresAt.Time),
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return &auth.RefreshTokenResponse{
			Success: false,
			Message: "refresh token is required",
		}, nil
	}

	tokens, err := s.tokenManager.RefreshTokens(req.RefreshToken)
	if err != nil {
		return &auth.RefreshTokenResponse{
			Success: false,
			Message: "invalid refresh token",
		}, nil
	}

	return &auth.RefreshTokenResponse{
		Success:      true,
		Message:      "tokens refreshed successfully",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    timestamppb.New(tokens.ExpiresAt),
	}, nil
}

// GetUserFromToken extracts user information from token
func (s *AuthService) GetUserFromToken(ctx context.Context, token string) (*entities.User, error) {
	claims, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// validateCredentials проверяет валидность логина и пароля
func validateCredentials(login, password string) error {
	if login == "" || password == "" {
		return errors.New("login and password are required")
	}
	return nil
}

// generateUserTokens генерирует токены для пользователя
func (s *AuthService) generateUserTokens(user *entities.User) (*TokenPair, error) {
	tokens, err := s.tokenManager.GenerateTokens(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	return tokens, nil
}

// convertUserToProto конвертирует доменного пользователя в protobuf
func convertUserToProto(user *entities.User) *auth.User {
	return &auth.User{
		Id:        user.ID,
		Login:     user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
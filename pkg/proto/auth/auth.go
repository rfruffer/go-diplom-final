package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// AuthServiceServer is the interface for auth service implementation
type AuthServiceServer interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
}

// AuthServiceClient is the interface for auth service client
type AuthServiceClient interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// RegisterResponse represents user registration response
type RegisterResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse represents user login response
type LoginResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *User     `json:"user"`
}

// ValidateTokenRequest represents token validation request
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// ValidateTokenResponse represents token validation response
type ValidateTokenResponse struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"user_id"`
	Login     string    `json:"login"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents token refresh response
type RefreshTokenResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// User represents user information
type User struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterAuthServiceServer регистрирует AuthService сервер в gRPC
func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	// TODO: заменить на сгенерированную protobuf реализацию
	// Пока что заглушка для компиляции
}
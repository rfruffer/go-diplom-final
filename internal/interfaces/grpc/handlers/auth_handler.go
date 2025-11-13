package handlers

import (
	"context"

	"github.com/fylgushev/go-diplom-final/internal/domain/services"
	"github.com/fylgushev/go-diplom-final/pkg/proto/auth"
)

// AuthHandler implements auth gRPC service
type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return h.authService.Register(ctx, req)
}

// Login handles user login
func (h *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return h.authService.Login(ctx, req)
}

// ValidateToken handles token validation
func (h *AuthHandler) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	return h.authService.ValidateToken(ctx, req)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	return h.authService.RefreshToken(ctx, req)
}
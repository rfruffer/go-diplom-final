package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/fylgushev/go-diplom-final/internal/domain/services"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserIDContextKey is the context key for user ID
	UserIDContextKey ContextKey = "user_id"
	// LoginContextKey is the context key for user login
	LoginContextKey ContextKey = "login"
)

// TokenManager interface for JWT operations
type TokenManager interface {
	ValidateToken(tokenString string) (*services.Claims, error)
}

// Claims is an alias for services.Claims
type Claims = services.Claims

// AuthInterceptor provides JWT authentication middleware for gRPC
type AuthInterceptor struct {
	tokenManager TokenManager
	exemptMethods map[string]bool
}

// NewAuthInterceptor creates a new authentication interceptor
func NewAuthInterceptor(tokenManager TokenManager) *AuthInterceptor {
	// Methods that don't require authentication
	exemptMethods := map[string]bool{
		"/auth.AuthService/Register":      true,
		"/auth.AuthService/Login":         true,
		"/auth.AuthService/RefreshToken":  true,
	}

	return &AuthInterceptor{
		tokenManager:  tokenManager,
		exemptMethods: exemptMethods,
	}
}

// UnaryInterceptor returns a gRPC unary server interceptor for JWT authentication
func (a *AuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for exempt methods
		if a.exemptMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract and validate token
		newCtx, err := a.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// StreamInterceptor returns a gRPC stream server interceptor for JWT authentication
func (a *AuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Skip authentication for exempt methods
		if a.exemptMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		// Extract and validate token
		newCtx, err := a.authenticate(ss.Context())
		if err != nil {
			return err
		}

		// Create new server stream with authenticated context
		wrapped := &authenticatedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		return handler(srv, wrapped)
	}
}

// authenticate extracts and validates JWT token from gRPC metadata
func (a *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	// Extract metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Get authorization header
	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	// Extract token from "Bearer <token>" format
	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	// Validate token
	claims, err := a.tokenManager.ValidateToken(token)
	if err != nil {
		if err.Error() == "token expired" {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Add user information to context
	newCtx := context.WithValue(ctx, UserIDContextKey, claims.UserID)
	newCtx = context.WithValue(newCtx, LoginContextKey, claims.Login)

	return newCtx, nil
}

// authenticatedServerStream wraps grpc.ServerStream with authenticated context
type authenticatedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the authenticated context
func (s *authenticatedServerStream) Context() context.Context {
	return s.ctx
}

// GetUserIDFromContext extracts user ID from authenticated context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

// GetLoginFromContext extracts user login from authenticated context
func GetLoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(LoginContextKey).(string)
	return login, ok
}

// AddExemptMethod adds a method to the exempt list (methods that don't require auth)
func (a *AuthInterceptor) AddExemptMethod(method string) {
	a.exemptMethods[method] = true
}

// RemoveExemptMethod removes a method from the exempt list
func (a *AuthInterceptor) RemoveExemptMethod(method string) {
	delete(a.exemptMethods, method)
}

// IsExemptMethod checks if a method is exempt from authentication
func (a *AuthInterceptor) IsExemptMethod(method string) bool {
	return a.exemptMethods[method]
}
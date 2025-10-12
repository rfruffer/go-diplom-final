package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// mockTokenManager implements TokenManager interface for testing
type mockTokenManager struct {
	validTokens map[string]*Claims
}

func (m *mockTokenManager) ValidateToken(token string) (*Claims, error) {
	if claims, ok := m.validTokens[token]; ok {
		return claims, nil
	}
	if token == "expired" {
		return nil, errors.New("token expired")
	}
	return nil, errors.New("invalid token")
}

func TestAuthInterceptor_UnaryInterceptor(t *testing.T) {
	userID := "user-123"
	login := "testuser"
	
	tm := &mockTokenManager{
		validTokens: map[string]*Claims{
			"valid-token": {
				UserID: userID,
				Login:  login,
			},
		},
	}
	interceptor := NewAuthInterceptor(tm)

	tests := []struct {
		name         string
		method       string
		metadata     metadata.MD
		wantErr      bool
		wantCode     codes.Code
		shouldSkip   bool
	}{
		{
			name:       "exempt method - register",
			method:     "/auth.AuthService/Register",
			metadata:   metadata.New(map[string]string{}),
			shouldSkip: true,
		},
		{
			name:       "exempt method - login",
			method:     "/auth.AuthService/Login",
			metadata:   metadata.New(map[string]string{}),
			shouldSkip: true,
		},
		{
			name:     "valid token",
			method:   "/data.DataService/StoreData",
			metadata: metadata.New(map[string]string{"authorization": "Bearer valid-token"}),
		},
		{
			name:     "missing metadata",
			method:   "/data.DataService/StoreData",
			metadata: nil,
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "missing authorization header",
			method:   "/data.DataService/StoreData",
			metadata: metadata.New(map[string]string{}),
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "invalid authorization format",
			method:   "/data.DataService/StoreData",
			metadata: metadata.New(map[string]string{"authorization": "InvalidFormat"}),
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "invalid token",
			method:   "/data.DataService/StoreData",
			metadata: metadata.New(map[string]string{"authorization": "Bearer invalid.token.here"}),
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create context with metadata
			ctx := context.Background()
			if tt.metadata != nil {
				ctx = metadata.NewIncomingContext(ctx, tt.metadata)
			}

			// Mock handler
			handlerCalled := false
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				handlerCalled = true

				if !tt.shouldSkip {
					// For non-exempt methods, check if user info is in context
					userIDFromCtx, hasUserID := GetUserIDFromContext(ctx)
					loginFromCtx, hasLogin := GetLoginFromContext(ctx)

					assert.True(t, hasUserID, "user ID should be in context")
					assert.True(t, hasLogin, "login should be in context")
					assert.Equal(t, userID, userIDFromCtx)
					assert.Equal(t, login, loginFromCtx)
				}

				return "success", nil
			}

			// Create interceptor
			unaryInterceptor := interceptor.UnaryInterceptor()

			// Create server info
			info := &grpc.UnaryServerInfo{
				FullMethod: tt.method,
			}

			// Call interceptor
			resp, err := unaryInterceptor(ctx, "request", info, handler)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantCode, status.Code(err))
				assert.Nil(t, resp)
				assert.False(t, handlerCalled)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "success", resp)
			assert.True(t, handlerCalled)
		})
	}
}

func TestAuthInterceptor_ExemptMethods(t *testing.T) {
	tm := &mockTokenManager{validTokens: make(map[string]*Claims)}
	interceptor := NewAuthInterceptor(tm)

	// Test default exempt methods
	assert.True(t, interceptor.IsExemptMethod("/auth.AuthService/Register"))
	assert.True(t, interceptor.IsExemptMethod("/auth.AuthService/Login"))
	assert.True(t, interceptor.IsExemptMethod("/auth.AuthService/RefreshToken"))
	assert.False(t, interceptor.IsExemptMethod("/data.DataService/StoreData"))

	// Test adding exempt method
	interceptor.AddExemptMethod("/custom.Service/Method")
	assert.True(t, interceptor.IsExemptMethod("/custom.Service/Method"))

	// Test removing exempt method
	interceptor.RemoveExemptMethod("/auth.AuthService/Register")
	assert.False(t, interceptor.IsExemptMethod("/auth.AuthService/Register"))
}

func TestGetUserInfoFromContext(t *testing.T) {
	ctx := context.Background()

	// Test empty context
	userID, hasUserID := GetUserIDFromContext(ctx)
	login, hasLogin := GetLoginFromContext(ctx)
	assert.False(t, hasUserID)
	assert.False(t, hasLogin)
	assert.Empty(t, userID)
	assert.Empty(t, login)

	// Test context with user info
	ctx = context.WithValue(ctx, UserIDContextKey, "user-123")
	ctx = context.WithValue(ctx, LoginContextKey, "testuser")

	userID, hasUserID = GetUserIDFromContext(ctx)
	login, hasLogin = GetLoginFromContext(ctx)
	assert.True(t, hasUserID)
	assert.True(t, hasLogin)
	assert.Equal(t, "user-123", userID)
	assert.Equal(t, "testuser", login)
}

func TestAuthInterceptor_ExpiredToken(t *testing.T) {
	tm := &mockTokenManager{validTokens: make(map[string]*Claims)}
	interceptor := NewAuthInterceptor(tm)

	// Create context with expired token
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer expired"}),
	)

	// Mock handler
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	// Create interceptor and server info
	unaryInterceptor := interceptor.UnaryInterceptor()
	info := &grpc.UnaryServerInfo{
		FullMethod: "/data.DataService/StoreData",
	}

	// Call interceptor with expired token
	resp, err := unaryInterceptor(ctx, "request", info, handler)

	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
	assert.Contains(t, err.Error(), "token expired")
	assert.Nil(t, resp)
}
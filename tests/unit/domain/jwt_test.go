package domain_test

import (
	"testing"
	"time"

	"github.com/fylgushev/go-diplom-final/internal/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenManager_GenerateTokens(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		login    string
		wantErr  bool
		errorMsg string
	}{
		{
			name:   "valid user",
			userID: "user-123",
			login:  "testuser",
		},
		{
			name:     "empty user ID",
			userID:   "",
			login:    "testuser",
			wantErr:  true,
			errorMsg: "invalid user ID",
		},
	}

	signingKey := []byte("test-secret-key-that-is-long-enough")
	tm := services.NewTokenManager(signingKey, 15*time.Minute, 24*time.Hour, "test-issuer")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tm.GenerateTokens(tt.userID, tt.login)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, tokens)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tokens)
			assert.NotEmpty(t, tokens.AccessToken)
			assert.NotEmpty(t, tokens.RefreshToken)
			assert.True(t, time.Now().Before(tokens.ExpiresAt))
		})
	}
}

func TestTokenManager_ValidateToken(t *testing.T) {
	signingKey := []byte("test-secret-key-that-is-long-enough")
	tm := services.NewTokenManager(signingKey, 15*time.Minute, 24*time.Hour, "test-issuer")

	userID := "user-123"
	login := "testuser"

	// Generate valid token
	tokens, err := tm.GenerateTokens(userID, login)
	require.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		wantErr   bool
		wantClaims *services.Claims
	}{
		{
			name:  "valid token",
			token: tokens.AccessToken,
			wantClaims: &services.Claims{
				UserID: userID,
				Login:  login,
			},
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := tm.ValidateToken(tt.token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, claims)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, claims)
			assert.Equal(t, tt.wantClaims.UserID, claims.UserID)
			assert.Equal(t, tt.wantClaims.Login, claims.Login)
		})
	}
}

func TestTokenManager_RefreshTokens(t *testing.T) {
	signingKey := []byte("test-secret-key-that-is-long-enough")
	tm := services.NewTokenManager(signingKey, 15*time.Minute, 24*time.Hour, "test-issuer")

	userID := "user-123"
	login := "testuser"

	// Generate initial tokens
	initialTokens, err := tm.GenerateTokens(userID, login)
	require.NoError(t, err)

	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
	}{
		{
			name:         "valid refresh token",
			refreshToken: initialTokens.RefreshToken,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid.token.here",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTokens, err := tm.RefreshTokens(tt.refreshToken)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, newTokens)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, newTokens)
			assert.NotEmpty(t, newTokens.AccessToken)
			assert.NotEmpty(t, newTokens.RefreshToken)
			
			// New tokens should be different from initial ones
			assert.NotEqual(t, initialTokens.AccessToken, newTokens.AccessToken)
			assert.NotEqual(t, initialTokens.RefreshToken, newTokens.RefreshToken)
		})
	}
}

func TestTokenManager_ExtractUserID(t *testing.T) {
	signingKey := []byte("test-secret-key-that-is-long-enough")
	tm := services.NewTokenManager(signingKey, 15*time.Minute, 24*time.Hour, "test-issuer")

	userID := "user-123"
	login := "testuser"

	tokens, err := tm.GenerateTokens(userID, login)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		wantUserID string
		wantErr    bool
	}{
		{
			name:       "valid token",
			token:      tokens.AccessToken,
			wantUserID: userID,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractedUserID, err := tm.ExtractUserID(tt.token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, extractedUserID)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUserID, extractedUserID)
		})
	}
}

func TestTokenManager_ExpiredToken(t *testing.T) {
	signingKey := []byte("test-secret-key-that-is-long-enough")
	// Create token manager with very short TTL
	tm := services.NewTokenManager(signingKey, 1*time.Millisecond, 24*time.Hour, "test-issuer")

	userID := "user-123"
	login := "testuser"

	tokens, err := tm.GenerateTokens(userID, login)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	_, err = tm.ValidateToken(tokens.AccessToken)
	assert.ErrorIs(t, err, services.ErrTokenExpired)
}

func TestDefaultConfig(t *testing.T) {
	config := services.DefaultConfig()

	assert.Equal(t, 15*time.Minute, config.AccessTokenTTL)
	assert.Equal(t, 24*time.Hour*7, config.RefreshTokenTTL)
	assert.Equal(t, "gophkeeper", config.Issuer)
}

func TestNewTokenManagerFromConfig(t *testing.T) {
	config := services.Config{
		SigningKey:      []byte("test-key"),
		AccessTokenTTL:  10 * time.Minute,
		RefreshTokenTTL: 48 * time.Hour,
		Issuer:          "test-app",
	}

	tm := services.NewTokenManagerFromConfig(config)

	assert.NotNil(t, tm)
}
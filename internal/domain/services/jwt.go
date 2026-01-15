package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidSigningKey = errors.New("invalid signing key")
	ErrInvalidUserID     = errors.New("invalid user ID")
)

// Claims represents JWT claims with custom fields
type Claims struct {
	UserID string `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

// TokenManager handles JWT token operations
type TokenManager struct {
	signingKey       []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	issuer           string
}

// NewTokenManager creates a new token manager
func NewTokenManager(signingKey []byte, accessTTL, refreshTTL time.Duration, issuer string) *TokenManager {
	return &TokenManager{
		signingKey:       signingKey,
		accessTokenTTL:   accessTTL,
		refreshTokenTTL:  refreshTTL,
		issuer:           issuer,
	}
}

// GenerateTokens generates both access and refresh tokens for a user
func (tm *TokenManager) GenerateTokens(userID, login string) (*TokenPair, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	now := time.Now()
	jti := uuid.New().String()

	// Generate access token
	accessToken, err := tm.generateToken(userID, login, jti, now, tm.accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (longer TTL, different JTI)
	refreshJTI := uuid.New().String()
	refreshToken, err := tm.generateToken(userID, login, refreshJTI, now, tm.refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(tm.accessTokenTTL),
	}, nil
}

// generateToken creates a JWT token with specified claims
func (tm *TokenManager) generateToken(userID, login, jti string, issuedAt time.Time, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			Issuer:    tm.issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(ttl)),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.signingKey)
}

// ValidateToken validates a JWT token and returns claims
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.signingKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshTokens generates new tokens using a valid refresh token
func (tm *TokenManager) RefreshTokens(refreshToken string) (*TokenPair, error) {
	claims, err := tm.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return tm.GenerateTokens(claims.UserID, claims.Login)
}

// ExtractUserID extracts user ID from a valid token
func (tm *TokenManager) ExtractUserID(tokenString string) (string, error) {
	claims, err := tm.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// DefaultConfig returns default token manager configuration
func DefaultConfig() Config {
	return Config{
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour * 7, // 7 days
		Issuer:          "gophkeeper",
	}
}

// Config represents token manager configuration
type Config struct {
	SigningKey      []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// NewTokenManagerFromConfig creates token manager from config
func NewTokenManagerFromConfig(config Config) *TokenManager {
	return NewTokenManager(
		config.SigningKey,
		config.AccessTokenTTL,
		config.RefreshTokenTTL,
		config.Issuer,
	)
}
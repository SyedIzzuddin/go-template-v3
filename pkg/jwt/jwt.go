package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
	ErrInvalidClaims = errors.New("invalid token claims")
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWTManager handles JWT token operations
type JWTManager struct {
	accessSecret     string
	refreshSecret    string
	accessExpiration time.Duration
	refreshExpiration time.Duration
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager(accessSecret, refreshSecret string, accessExp, refreshExp time.Duration) *JWTManager {
	return &JWTManager{
		accessSecret:      accessSecret,
		refreshSecret:     refreshSecret,
		accessExpiration:  accessExp,
		refreshExpiration: refreshExp,
	}
}

// GenerateTokenPair creates both access and refresh tokens for a user
func (jm *JWTManager) GenerateTokenPair(userID int, email string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := jm.generateToken(userID, email, jm.accessSecret, jm.accessExpiration)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := jm.generateToken(userID, email, jm.refreshSecret, jm.refreshExpiration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateToken creates a JWT token with the given parameters
func (jm *JWTManager) generateToken(userID int, email, secret string, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "go-template",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken validates an access token and returns the claims
func (jm *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return jm.validateToken(tokenString, jm.accessSecret)
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (jm *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return jm.validateToken(tokenString, jm.refreshSecret)
}

// validateToken validates a JWT token with the given secret
func (jm *JWTManager) validateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token using a valid refresh token
func (jm *JWTManager) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := jm.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Generate new access token with same user info
	return jm.generateToken(claims.UserID, claims.Email, jm.accessSecret, jm.accessExpiration)
}
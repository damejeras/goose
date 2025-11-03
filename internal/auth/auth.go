package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrMissingToken     = errors.New("missing token")
)

type Config struct {
	GoogleClientID string
	JWTSecret      []byte
	JWTExpiration  time.Duration
}

type Service struct {
	config Config
}

func NewService(config Config) *Service {
	if config.JWTExpiration == 0 {
		config.JWTExpiration = 24 * time.Hour // default 24 hours
	}
	return &Service{config: config}
}

// GoogleTokenInfo contains the validated information from a Google ID token
type GoogleTokenInfo struct {
	Email      string
	GoogleID   string
	Name       string
	Picture    string
	Verified   bool
}

// ValidateGoogleIDToken validates a Google ID token and extracts user information
func (s *Service) ValidateGoogleIDToken(ctx context.Context, idToken string) (*GoogleTokenInfo, error) {
	// This validates:
	// 1. Token signature (signed by Google)
	// 2. Token expiration
	// 3. Audience claim (aud) matches our Client ID
	payload, err := idtoken.Validate(ctx, idToken, s.config.GoogleClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Google ID token: %w", err)
	}

	// Additional validation: verify issuer is Google
	issuer, _ := payload.Claims["iss"].(string)
	if issuer != "https://accounts.google.com" && issuer != "accounts.google.com" {
		return nil, fmt.Errorf("invalid issuer: %s", issuer)
	}

	// Verify audience explicitly (belt and suspenders)
	audience, _ := payload.Claims["aud"].(string)
	if audience != s.config.GoogleClientID {
		return nil, fmt.Errorf("invalid audience: expected %s, got %s", s.config.GoogleClientID, audience)
	}

	email, _ := payload.Claims["email"].(string)
	sub, _ := payload.Claims["sub"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)

	if email == "" || sub == "" {
		return nil, ErrInvalidToken
	}

	return &GoogleTokenInfo{
		Email:      email,
		GoogleID:   sub,
		Name:       name,
		Picture:    picture,
		Verified:   emailVerified,
	}, nil
}

// JWTClaims represents the JWT claims for our application
type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT token for a user
func (s *Service) GenerateJWT(userID int64, email string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.config.JWTSecret)
}

// ValidateJWT validates a JWT token and returns the claims
func (s *Service) ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.JWTSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GenerateRandomSecret generates a random secret for JWT signing (useful for development)
func GenerateRandomSecret() ([]byte, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// GenerateRandomSecretString generates a random secret as a base64 string
func GenerateRandomSecretString() (string, error) {
	secret, err := GenerateRandomSecret()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(secret), nil
}

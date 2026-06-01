package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/model"
)

const minPasswordLength = 8

// userStore is the slice of UserRepository that AuthService needs. Declaring
// it here keeps the service testable with a fake.
type userStore interface {
	Create(ctx context.Context, email, passwordHash string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, string, error)
	GetByID(ctx context.Context, id uint64) (model.User, error)
}

// refreshStore persists refresh-token hashes. The plaintext token never
// touches the database (only its SHA-256 hash).
type refreshStore interface {
	Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error
	FindValidUser(ctx context.Context, tokenHash string) (uint64, error)
	Revoke(ctx context.Context, tokenHash string) error
}

// AuthService implements registration, login, token refresh and logout.
type AuthService struct {
	users      userStore
	refresh    refreshStore
	tokens     *auth.TokenManager
	refreshTTL time.Duration
}

// NewAuthService wires an AuthService to its stores, token manager and the
// refresh-token lifetime.
func NewAuthService(users userStore, refresh refreshStore, tokens *auth.TokenManager, refreshTTL time.Duration) *AuthService {
	return &AuthService{users: users, refresh: refresh, tokens: tokens, refreshTTL: refreshTTL}
}

// Register creates an account and returns a short-lived access token (in the
// AuthResponse) plus a long-lived refresh token (the second return value, which
// the handler stores in an HttpOnly cookie).
func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (model.AuthResponse, string, error) {
	if req.Email == "" {
		return model.AuthResponse{}, "", apperror.Validation("email is required")
	}
	if len(req.Password) < minPasswordLength {
		return model.AuthResponse{}, "", apperror.Validation("password must be at least 8 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	user, err := s.users.Create(ctx, req.Email, string(hash))
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	return s.issue(ctx, user)
}

// Login verifies credentials and returns an access token plus a refresh token.
// Invalid credentials return apperror.ErrUnauthorized.
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, string, error) {
	if req.Email == "" || req.Password == "" {
		return model.AuthResponse{}, "", apperror.Validation("email and password are required")
	}
	user, hash, err := s.users.GetByEmail(ctx, req.Email)
	if errors.Is(err, apperror.ErrNotFound) {
		return model.AuthResponse{}, "", apperror.ErrUnauthorized
	}
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		return model.AuthResponse{}, "", apperror.ErrUnauthorized
	}
	return s.issue(ctx, user)
}

// Refresh validates a refresh token, rotates it (the old one is revoked and a
// new one issued) and returns a fresh access token. An invalid, expired or
// revoked token returns apperror.ErrUnauthorized.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (model.AuthResponse, string, error) {
	if refreshToken == "" {
		return model.AuthResponse{}, "", apperror.ErrUnauthorized
	}
	hash := hashToken(refreshToken)
	userID, err := s.refresh.FindValidUser(ctx, hash)
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	// Token rotation: the used refresh token is immediately revoked so it can
	// never be replayed.
	if err := s.refresh.Revoke(ctx, hash); err != nil {
		return model.AuthResponse{}, "", err
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	return s.issue(ctx, user)
}

// Logout revokes the given refresh token. It is idempotent: an unknown token
// is not an error.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.refresh.Revoke(ctx, hashToken(refreshToken))
}

// issue creates an access token and a new refresh token for the user.
func (s *AuthService) issue(ctx context.Context, user model.User) (model.AuthResponse, string, error) {
	return issueTokens(ctx, s.tokens, s.refresh, s.refreshTTL, user)
}

// issueTokens creates an access token and a new (stored, hashed) refresh token
// for the user. It is shared by AuthService and OAuthService. It returns the
// AuthResponse and the plaintext refresh token (never stored).
func issueTokens(ctx context.Context, tokens *auth.TokenManager, refresh refreshStore, refreshTTL time.Duration, user model.User) (model.AuthResponse, string, error) {
	access, err := tokens.Issue(user.ID, time.Now())
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	plaintext, hash, err := newRefreshToken()
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	if err := refresh.Create(ctx, user.ID, hash, time.Now().Add(refreshTTL)); err != nil {
		return model.AuthResponse{}, "", err
	}
	return model.AuthResponse{Token: access, User: user}, plaintext, nil
}

// newRefreshToken returns a random 64-char hex token and its SHA-256 hash.
func newRefreshToken() (plaintext, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	plaintext = hex.EncodeToString(buf)
	return plaintext, hashToken(plaintext), nil
}

// hashToken returns the SHA-256 hex digest of a token. SHA-256 (not bcrypt) is
// fine here because the token is already high-entropy random.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

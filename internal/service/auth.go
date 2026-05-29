package service

import (
	"context"
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
}

// AuthService implements registration and login.
type AuthService struct {
	users  userStore
	tokens *auth.TokenManager
}

// NewAuthService wires an AuthService to its user store and token manager.
func NewAuthService(users userStore, tokens *auth.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

// Register validates the request, hashes the password, stores the user and
// returns a signed token.
func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (model.AuthResponse, error) {
	if req.Email == "" {
		return model.AuthResponse{}, apperror.Validation("email is required")
	}
	if len(req.Password) < minPasswordLength {
		return model.AuthResponse{}, apperror.Validation("password must be at least 8 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthResponse{}, err
	}
	user, err := s.users.Create(ctx, req.Email, string(hash))
	if err != nil {
		return model.AuthResponse{}, err
	}
	return s.issue(user)
}

// Login verifies credentials and returns a signed token. Invalid credentials
// (unknown email or wrong password) return apperror.ErrUnauthorized.
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return model.AuthResponse{}, apperror.Validation("email and password are required")
	}
	user, hash, err := s.users.GetByEmail(ctx, req.Email)
	if errors.Is(err, apperror.ErrNotFound) {
		return model.AuthResponse{}, apperror.ErrUnauthorized
	}
	if err != nil {
		return model.AuthResponse{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		return model.AuthResponse{}, apperror.ErrUnauthorized
	}
	return s.issue(user)
}

func (s *AuthService) issue(user model.User) (model.AuthResponse, error) {
	token, err := s.tokens.Issue(user.ID, time.Now())
	if err != nil {
		return model.AuthResponse{}, err
	}
	return model.AuthResponse{Token: token, User: user}, nil
}

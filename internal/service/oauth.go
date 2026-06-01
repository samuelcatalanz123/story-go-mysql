package service

import (
	"context"
	"errors"
	"time"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/model"
)

// oauthUserStore is the subset of UserRepository that OAuthService needs.
type oauthUserStore interface {
	GetByID(ctx context.Context, id uint64) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, string, error)
	CreateOAuthUser(ctx context.Context, email string) (model.User, error)
}

// oauthAccountStore links users to external providers.
type oauthAccountStore interface {
	FindUserID(ctx context.Context, provider, subject string) (uint64, error)
	Link(ctx context.Context, userID uint64, provider, subject string) error
}

// googleAuthenticator verifies a Google authorization code and returns the
// user's immutable subject and email. Implemented by oauth.GoogleAuthenticator.
type googleAuthenticator interface {
	ProviderName() string
	Exchange(ctx context.Context, code, codeVerifier string) (subject, email string, err error)
}

// OAuthService handles "Sign in with Google". google may be nil when no Google
// credentials are configured, in which case LoginWithGoogle reports that.
type OAuthService struct {
	users      oauthUserStore
	accounts   oauthAccountStore
	google     googleAuthenticator
	tokens     *auth.TokenManager
	refresh    refreshStore
	refreshTTL time.Duration
}

// NewOAuthService wires an OAuthService. Pass google=nil to disable Google login.
func NewOAuthService(users oauthUserStore, accounts oauthAccountStore, google googleAuthenticator, tokens *auth.TokenManager, refresh refreshStore, refreshTTL time.Duration) *OAuthService {
	return &OAuthService{users: users, accounts: accounts, google: google, tokens: tokens, refresh: refresh, refreshTTL: refreshTTL}
}

// LoginWithGoogle verifies the Google code, finds or creates the matching user
// and issues our own access + refresh tokens.
func (s *OAuthService) LoginWithGoogle(ctx context.Context, code, codeVerifier string) (model.AuthResponse, string, error) {
	if s.google == nil {
		return model.AuthResponse{}, "", apperror.Validation("google login is not configured on the server")
	}
	if code == "" || codeVerifier == "" {
		return model.AuthResponse{}, "", apperror.Validation("code and codeVerifier are required")
	}

	subject, email, err := s.google.Exchange(ctx, code, codeVerifier)
	if err != nil {
		return model.AuthResponse{}, "", apperror.ErrUnauthorized
	}

	user, err := s.resolveUser(ctx, s.google.ProviderName(), subject, email)
	if err != nil {
		return model.AuthResponse{}, "", err
	}
	return issueTokens(ctx, s.tokens, s.refresh, s.refreshTTL, user)
}

// resolveUser implements the login-or-signup rules:
//  1. (provider, subject) already linked    -> log that user in.
//  2. email exists but not linked           -> link it and log in.
//  3. nothing exists                         -> create user + link.
func (s *OAuthService) resolveUser(ctx context.Context, provider, subject, email string) (model.User, error) {
	userID, err := s.accounts.FindUserID(ctx, provider, subject)
	if err == nil {
		return s.users.GetByID(ctx, userID)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		return model.User{}, err
	}

	// Not linked yet: reuse an existing account with the same email, or create one.
	user, _, err := s.users.GetByEmail(ctx, email)
	if errors.Is(err, apperror.ErrNotFound) {
		user, err = s.users.CreateOAuthUser(ctx, email)
	}
	if err != nil {
		return model.User{}, err
	}

	if err := s.accounts.Link(ctx, user.ID, provider, subject); err != nil {
		return model.User{}, err
	}
	return user, nil
}

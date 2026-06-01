package service

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"story-go-mysql/internal/email"
	"story-go-mysql/internal/model"
)

// emailVerifyUserStore is the subset of UserRepository this service needs.
type emailVerifyUserStore interface {
	GetByID(ctx context.Context, id uint64) (model.User, error)
	SetEmailVerified(ctx context.Context, id uint64) error
}

// emailVerifyTokenStore stores and consumes one-time verification tokens.
type emailVerifyTokenStore interface {
	Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error
	Consume(ctx context.Context, tokenHash string) (uint64, error)
}

var verifyEmailTmpl = template.Must(template.New("verify").Parse(
	`<p>¡Bienvenido!</p>
<p>Confirma tu correo haciendo clic en el enlace:</p>
<p><a href="{{.Link}}">Verificar mi correo</a></p>
<p>El enlace caduca en 24 horas.</p>`))

// EmailVerificationService implements the "verify your email" flow.
type EmailVerificationService struct {
	users      emailVerifyUserStore
	tokens     emailVerifyTokenStore
	mailer     email.Sender
	appBaseURL string
	ttl        time.Duration
}

// NewEmailVerificationService wires the service to its dependencies.
func NewEmailVerificationService(users emailVerifyUserStore, tokens emailVerifyTokenStore, mailer email.Sender, appBaseURL string, ttl time.Duration) *EmailVerificationService {
	return &EmailVerificationService{users: users, tokens: tokens, mailer: mailer, appBaseURL: appBaseURL, ttl: ttl}
}

// Send generates a verification token for the user and emails them the link.
// It looks up the user's email itself, so callers only need the user ID.
func (s *EmailVerificationService) Send(ctx context.Context, userID uint64) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	plaintext, hash, err := newRefreshToken() // random token + SHA-256 hash
	if err != nil {
		return err
	}
	if err := s.tokens.Create(ctx, userID, hash, time.Now().Add(s.ttl)); err != nil {
		return err
	}
	var body bytes.Buffer
	link := s.appBaseURL + "/verify-email?token=" + plaintext
	if err := verifyEmailTmpl.Execute(&body, struct{ Link string }{Link: link}); err != nil {
		return err
	}
	return s.mailer.Send(ctx, email.Message{
		To:      user.Email,
		Subject: "Verifica tu correo",
		HTML:    body.String(),
	})
}

// Verify consumes the token (single use) and marks the user's email verified.
func (s *EmailVerificationService) Verify(ctx context.Context, token string) error {
	userID, err := s.tokens.Consume(ctx, hashToken(token))
	if err != nil {
		return err
	}
	return s.users.SetEmailVerified(ctx, userID)
}

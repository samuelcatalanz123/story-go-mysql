package service

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"golang.org/x/crypto/bcrypt"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/email"
	"story-go-mysql/internal/model"
)

// passwordResetUserStore is the subset of UserRepository this service needs.
type passwordResetUserStore interface {
	GetByEmail(ctx context.Context, email string) (model.User, string, error)
	UpdatePassword(ctx context.Context, id uint64, passwordHash string) error
}

// passwordResetTokenStore stores and consumes one-time reset tokens.
type passwordResetTokenStore interface {
	Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error
	Consume(ctx context.Context, tokenHash string) (uint64, error)
}

// resetEmailTmpl is the HTML body of the reset email. html/template escapes
// the data automatically, preventing injection into the markup.
var resetEmailTmpl = template.Must(template.New("reset").Parse(
	`<p>Hola,</p>
<p>Recibimos una solicitud para restablecer tu contraseña. Haz clic en el enlace:</p>
<p><a href="{{.Link}}">Restablecer mi contraseña</a></p>
<p>El enlace caduca en una hora. Si no fuiste tú, ignora este correo.</p>`))

// PasswordResetService implements the forgot/reset password flow.
type PasswordResetService struct {
	users      passwordResetUserStore
	tokens     passwordResetTokenStore
	mailer     email.Sender
	appBaseURL string
	ttl        time.Duration
}

// NewPasswordResetService wires the service to its dependencies.
func NewPasswordResetService(users passwordResetUserStore, tokens passwordResetTokenStore, mailer email.Sender, appBaseURL string, ttl time.Duration) *PasswordResetService {
	return &PasswordResetService{users: users, tokens: tokens, mailer: mailer, appBaseURL: appBaseURL, ttl: ttl}
}

// ForgotPassword sends a reset link if the email belongs to a user. It always
// reports success (returns nil) regardless of whether the email exists, so an
// attacker can't use this endpoint to discover which emails are registered.
func (s *PasswordResetService) ForgotPassword(ctx context.Context, emailAddr string) error {
	user, _, err := s.users.GetByEmail(ctx, emailAddr)
	if errors.Is(err, apperror.ErrNotFound) {
		return nil // email not registered: pretend success
	}
	if err != nil {
		return err
	}

	plaintext, hash, err := newRefreshToken() // random token + its SHA-256 hash
	if err != nil {
		return err
	}
	if err := s.tokens.Create(ctx, user.ID, hash, time.Now().Add(s.ttl)); err != nil {
		return err
	}

	var body bytes.Buffer
	link := s.appBaseURL + "/reset-password?token=" + plaintext
	if err := resetEmailTmpl.Execute(&body, struct{ Link string }{Link: link}); err != nil {
		return err
	}
	return s.mailer.Send(ctx, email.Message{
		To:      user.Email,
		Subject: "Restablece tu contraseña",
		HTML:    body.String(),
	})
}

// ResetPassword consumes the token (single use) and sets the new password.
// An invalid, expired or already-used token returns apperror.ErrUnauthorized.
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < minPasswordLength {
		return apperror.Validation("password must be at least 8 characters")
	}
	userID, err := s.tokens.Consume(ctx, hashToken(token))
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(ctx, userID, string(hash))
}

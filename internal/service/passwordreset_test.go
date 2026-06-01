package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/email"
	"story-go-mysql/internal/model"
)

type fakeResetUserStore struct {
	users   map[string]model.User
	updated map[uint64]string
}

func (f *fakeResetUserStore) GetByEmail(_ context.Context, e string) (model.User, string, error) {
	u, ok := f.users[e]
	if !ok {
		return model.User{}, "", apperror.ErrNotFound
	}
	return u, "", nil
}

func (f *fakeResetUserStore) UpdatePassword(_ context.Context, id uint64, hash string) error {
	f.updated[id] = hash
	return nil
}

type fakeResetTokenStore struct {
	byHash map[string]uint64
	used   map[string]bool
}

func (f *fakeResetTokenStore) Create(_ context.Context, userID uint64, hash string, _ time.Time) error {
	f.byHash[hash] = userID
	return nil
}

func (f *fakeResetTokenStore) Consume(_ context.Context, hash string) (uint64, error) {
	userID, ok := f.byHash[hash]
	if !ok || f.used[hash] {
		return 0, apperror.ErrUnauthorized
	}
	f.used[hash] = true
	return userID, nil
}

func newResetService(mailer *fakeMailer) *PasswordResetService {
	users := &fakeResetUserStore{
		users:   map[string]model.User{"a@b.com": {ID: 1, Email: "a@b.com"}},
		updated: map[uint64]string{},
	}
	tokens := &fakeResetTokenStore{byHash: map[string]uint64{}, used: map[string]bool{}}
	return NewPasswordResetService(users, tokens, mailer, "http://localhost:8090", time.Hour)
}

type fakeMailer struct {
	last  email.Message
	count int
}

func (m *fakeMailer) Send(_ context.Context, msg email.Message) error {
	m.last = msg
	m.count++
	return nil
}

// linkToken extrae el token del enlace dentro del HTML del correo.
func linkToken(html string) string {
	const marker = "token="
	i := strings.Index(html, marker)
	if i < 0 {
		return ""
	}
	rest := html[i+len(marker):]
	if end := strings.IndexByte(rest, '"'); end >= 0 {
		return rest[:end]
	}
	return rest
}

func TestForgotPasswordUnknownEmailStillSucceeds(t *testing.T) {
	mailer := &fakeMailer{}
	svc := newResetService(mailer)
	if err := svc.ForgotPassword(context.Background(), "noexiste@x.com"); err != nil {
		t.Fatalf("debe devolver nil aunque el email no exista, obtuve %v", err)
	}
	if mailer.count != 0 {
		t.Fatal("no debería enviar correo a un email inexistente")
	}
}

func TestForgotAndResetPasswordFlow(t *testing.T) {
	mailer := &fakeMailer{}
	svc := newResetService(mailer)

	if err := svc.ForgotPassword(context.Background(), "a@b.com"); err != nil {
		t.Fatal(err)
	}
	if mailer.count != 1 {
		t.Fatal("esperaba que se enviara el correo de reset")
	}
	token := linkToken(mailer.last.HTML)
	if token == "" {
		t.Fatal("no encontré el token en el correo")
	}

	if err := svc.ResetPassword(context.Background(), token, "nuevaclave123"); err != nil {
		t.Fatalf("el reset debería funcionar, obtuve %v", err)
	}

	// Un segundo uso del mismo token debe fallar (un solo uso).
	if err := svc.ResetPassword(context.Background(), token, "otraclave123"); !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("el token ya usado no debería servir, obtuve %v", err)
	}
}

func TestResetPasswordRejectsShortPassword(t *testing.T) {
	svc := newResetService(&fakeMailer{})
	if err := svc.ResetPassword(context.Background(), "cualquiercosa", "corta"); err == nil {
		t.Fatal("esperaba error por contraseña corta")
	}
}

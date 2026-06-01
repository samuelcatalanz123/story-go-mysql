package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

type fakeVerifyUserStore struct {
	users    map[uint64]model.User
	verified map[uint64]bool
}

func (f *fakeVerifyUserStore) GetByID(_ context.Context, id uint64) (model.User, error) {
	u, ok := f.users[id]
	if !ok {
		return model.User{}, apperror.ErrNotFound
	}
	return u, nil
}

func (f *fakeVerifyUserStore) SetEmailVerified(_ context.Context, id uint64) error {
	f.verified[id] = true
	return nil
}

func TestEmailVerificationFlow(t *testing.T) {
	mailer := &fakeMailer{}
	users := &fakeVerifyUserStore{
		users:    map[uint64]model.User{1: {ID: 1, Email: "a@b.com"}},
		verified: map[uint64]bool{},
	}
	// fakeResetTokenStore (de passwordreset_test) implementa la misma interfaz.
	tokens := &fakeResetTokenStore{byHash: map[string]uint64{}, used: map[string]bool{}}
	svc := NewEmailVerificationService(users, tokens, mailer, "http://localhost:8090", time.Hour)

	if err := svc.Send(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	token := linkToken(mailer.last.HTML)
	if token == "" {
		t.Fatal("no encontré el token en el correo de verificación")
	}

	if err := svc.Verify(context.Background(), token); err != nil {
		t.Fatalf("la verificación debería funcionar, obtuve %v", err)
	}
	if !users.verified[1] {
		t.Fatal("el email del usuario debería quedar marcado como verificado")
	}

	// El token es de un solo uso.
	if err := svc.Verify(context.Background(), token); !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("el token ya usado no debería servir, obtuve %v", err)
	}
}

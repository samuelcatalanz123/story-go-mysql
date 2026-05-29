package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/model"
)

type fakeUserStore struct {
	users  map[string]model.User
	hashes map[string]string
	nextID uint64
}

func newFakeStore() *fakeUserStore {
	return &fakeUserStore{users: map[string]model.User{}, hashes: map[string]string{}}
}

func (f *fakeUserStore) Create(_ context.Context, email, hash string) (model.User, error) {
	if _, ok := f.users[email]; ok {
		return model.User{}, apperror.ErrDuplicateEmail
	}
	f.nextID++
	u := model.User{ID: f.nextID, Email: email}
	f.users[email] = u
	f.hashes[email] = hash
	return u, nil
}

func (f *fakeUserStore) GetByEmail(_ context.Context, email string) (model.User, string, error) {
	u, ok := f.users[email]
	if !ok {
		return model.User{}, "", apperror.ErrNotFound
	}
	return u, f.hashes[email], nil
}

func newAuthService() *AuthService {
	return NewAuthService(newFakeStore(), auth.NewTokenManager("secret", time.Hour))
}

func TestRegisterIssuesToken(t *testing.T) {
	s := newAuthService()
	res, err := s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "password123"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Token == "" {
		t.Fatal("esperaba un token")
	}
	if res.User.Email != "a@b.com" {
		t.Fatalf("email inesperado: %s", res.User.Email)
	}
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	s := newAuthService()
	_, err := s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "short"})
	var v apperror.ValidationError
	if !errors.As(err, &v) {
		t.Fatalf("esperaba ValidationError, obtuve %v", err)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	s := newAuthService()
	_, _ = s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "password123"})
	_, err := s.Login(context.Background(), model.LoginRequest{Email: "a@b.com", Password: "wrongpass1"})
	if !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("esperaba Unauthorized, obtuve %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	s := newAuthService()
	_, _ = s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "password123"})
	res, err := s.Login(context.Background(), model.LoginRequest{Email: "a@b.com", Password: "password123"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Token == "" {
		t.Fatal("esperaba un token")
	}
}

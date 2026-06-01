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

func (f *fakeUserStore) GetByID(_ context.Context, id uint64) (model.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return model.User{}, apperror.ErrNotFound
}

// fakeRefreshStore is an in-memory refreshStore for testing.
type fakeRefreshStore struct {
	userByHash  map[string]uint64
	revoked     map[string]bool
	expiredHash map[string]bool
}

func newFakeRefreshStore() *fakeRefreshStore {
	return &fakeRefreshStore{
		userByHash:  map[string]uint64{},
		revoked:     map[string]bool{},
		expiredHash: map[string]bool{},
	}
}

func (f *fakeRefreshStore) Create(_ context.Context, userID uint64, tokenHash string, _ time.Time) error {
	f.userByHash[tokenHash] = userID
	return nil
}

func (f *fakeRefreshStore) FindValidUser(_ context.Context, tokenHash string) (uint64, error) {
	userID, ok := f.userByHash[tokenHash]
	if !ok || f.revoked[tokenHash] || f.expiredHash[tokenHash] {
		return 0, apperror.ErrUnauthorized
	}
	return userID, nil
}

func (f *fakeRefreshStore) Revoke(_ context.Context, tokenHash string) error {
	f.revoked[tokenHash] = true
	return nil
}

func newAuthService() *AuthService {
	return NewAuthService(newFakeStore(), newFakeRefreshStore(), auth.NewTokenManager("secret", time.Hour), time.Hour)
}

func TestRegisterIssuesTokens(t *testing.T) {
	s := newAuthService()
	res, refresh, err := s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "password123"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Token == "" || refresh == "" {
		t.Fatal("esperaba access token y refresh token")
	}
	if res.User.Email != "a@b.com" {
		t.Fatalf("email inesperado: %s", res.User.Email)
	}
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	s := newAuthService()
	_, _, err := s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "short"})
	var v apperror.ValidationError
	if !errors.As(err, &v) {
		t.Fatalf("esperaba ValidationError, obtuve %v", err)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	s := newAuthService()
	_, _, _ = s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "password123"})
	_, _, err := s.Login(context.Background(), model.LoginRequest{Email: "a@b.com", Password: "wrongpass1"})
	if !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("esperaba Unauthorized, obtuve %v", err)
	}
}

func TestRefreshRotatesToken(t *testing.T) {
	s := newAuthService()
	_, refresh, _ := s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "password123"})

	// Usar el refresh token devuelve uno nuevo (rotación)...
	res2, refresh2, err := s.Refresh(context.Background(), refresh)
	if err != nil {
		t.Fatal(err)
	}
	if res2.Token == "" || refresh2 == "" || refresh2 == refresh {
		t.Fatal("esperaba un access y un refresh nuevos (distintos)")
	}

	// ...y el viejo ya no sirve (fue revocado al rotar).
	if _, _, err := s.Refresh(context.Background(), refresh); !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("el refresh viejo debería estar revocado, obtuve %v", err)
	}
}

func TestLogoutRevokesRefresh(t *testing.T) {
	s := newAuthService()
	_, refresh, _ := s.Register(context.Background(), model.RegisterRequest{Email: "a@b.com", Password: "password123"})

	if err := s.Logout(context.Background(), refresh); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Refresh(context.Background(), refresh); !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("tras logout el refresh no debería servir, obtuve %v", err)
	}
}

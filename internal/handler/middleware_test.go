package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/handler"
)

func TestRequireAuthRejectsMissingToken(t *testing.T) {
	tm := auth.NewTokenManager("secret", time.Hour)
	h := handler.RequireAuth(tm, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/characters", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperaba 401, obtuve %d", rec.Code)
	}
}

func TestRequireAuthAllowsValidToken(t *testing.T) {
	tm := auth.NewTokenManager("secret", time.Hour)
	tok, _ := tm.Issue(1, time.Now())
	called := false
	h := handler.RequireAuth(tm, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodPost, "/characters", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !called {
		t.Fatalf("esperaba 200 y next llamado; code=%d called=%v", rec.Code, called)
	}
}

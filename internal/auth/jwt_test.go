package auth

import (
	"testing"
	"time"
)

func TestIssueAndParse(t *testing.T) {
	m := NewTokenManager("secret", time.Hour)
	tok, err := m.Issue(42, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	id, err := m.Parse(tok)
	if err != nil {
		t.Fatal(err)
	}
	if id != 42 {
		t.Fatalf("esperaba 42, obtuve %d", id)
	}
}

func TestParseRejectsExpired(t *testing.T) {
	m := NewTokenManager("secret", time.Hour)
	tok, _ := m.Issue(1, time.Now().Add(-2*time.Hour)) // emitido hace 2h, ttl 1h
	if _, err := m.Parse(tok); err == nil {
		t.Fatal("esperaba error por token expirado")
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	issuer := NewTokenManager("secret-a", time.Hour)
	verifier := NewTokenManager("secret-b", time.Hour)
	tok, _ := issuer.Issue(1, time.Now())
	if _, err := verifier.Parse(tok); err == nil {
		t.Fatal("esperaba error por firma inválida")
	}
}

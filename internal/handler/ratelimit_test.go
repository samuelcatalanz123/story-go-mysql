package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"story-go-mysql/internal/cache"
)

func TestRateLimitBlocksAfterLimit(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	c, err := cache.NewRedis(context.Background(), mr.Addr(), "")
	if err != nil {
		t.Fatal(err)
	}

	// Permite 2 por ventana; la 3a debe ser bloqueada.
	limited := RateLimit(c, 2, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	codes := []int{}
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		rec := httptest.NewRecorder()
		limited.ServeHTTP(rec, req)
		codes = append(codes, rec.Code)
	}

	if codes[0] != 200 || codes[1] != 200 {
		t.Fatalf("las primeras 2 deberían pasar (200), obtuve %v", codes)
	}
	if codes[2] != http.StatusTooManyRequests {
		t.Fatalf("la 3a debería ser 429, obtuve %d", codes[2])
	}
}

func TestRateLimitSeparatesByIP(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	c, _ := cache.NewRedis(context.Background(), mr.Addr(), "")

	limited := RateLimit(c, 1, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Dos IPs distintas: cada una tiene su propio contador.
	for _, ip := range []string{"1.1.1.1:1", "2.2.2.2:2"} {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		limited.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("la primera petición de %s debería pasar, obtuve %d", ip, rec.Code)
		}
	}
}

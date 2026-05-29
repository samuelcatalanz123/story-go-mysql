package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"story-go-mysql/internal/handler"
)

func TestSPAHandlerServesExistingFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("JS"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := handler.SPAHandler(dir)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/app.js", nil))
	if rec.Body.String() != "JS" {
		t.Fatalf("esperaba el contenido del archivo, obtuve %q", rec.Body.String())
	}
}

func TestSPAHandlerFallsBackToIndex(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := handler.SPAHandler(dir)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/scenes", nil))
	if rec.Body.String() != "INDEX" {
		t.Fatalf("esperaba el index.html como fallback, obtuve %q", rec.Body.String())
	}
}

func TestWithFrontendRoutesAPI(t *testing.T) {
	api := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/characters" {
			_, _ = w.Write([]byte("API"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := handler.WithFrontend(api, dir)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/characters", nil))
	if rec.Body.String() != "API" {
		t.Fatalf("esperaba que /api/characters enrutara a la API, obtuve %q", rec.Body.String())
	}
}

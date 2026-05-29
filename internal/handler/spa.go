package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

// SPAHandler serves static files from dir. Requests that don't map to an
// existing file fall back to index.html, so client-side routing (React
// Router) keeps working when the user reloads a deep link.
func SPAHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	index := filepath.Join(dir, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, index)
	})
}

// WithFrontend composes a top-level handler: requests under /api go to the
// API (with the /api prefix stripped); everything else is served by the SPA
// handler when webDir exists. If webDir is empty or missing, only the API is
// served — useful in local development, where Vite serves the frontend.
func WithFrontend(api http.Handler, webDir string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))
	if webDir != "" {
		if info, err := os.Stat(webDir); err == nil && info.IsDir() {
			mux.Handle("/", SPAHandler(webDir))
		}
	}
	return mux
}

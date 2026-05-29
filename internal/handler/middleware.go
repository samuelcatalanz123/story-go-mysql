package handler

import (
	"net/http"
	"strings"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/web"
)

// RequireAuth wraps next, rejecting requests that lack a valid Bearer token.
func RequireAuth(tokens *auth.TokenManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token == "" {
			web.RespondError(w, authResource, apperror.ErrUnauthorized)
			return
		}
		if _, err := tokens.Parse(token); err != nil {
			web.RespondError(w, authResource, apperror.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

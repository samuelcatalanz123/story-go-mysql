package handler

import (
	"context"
	"net/http"
	"strings"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/web"
)

// contextKey is an unexported type for context keys, so values stored by this
// package can never collide with keys from other packages.
type contextKey string

const userIDKey contextKey = "userID"

// RequireAuth wraps next, rejecting requests that lack a valid Bearer token.
// On success it stores the authenticated user ID in the request context so
// downstream handlers can read it with userIDFrom.
func RequireAuth(tokens *auth.TokenManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token == "" {
			web.RespondError(w, authResource, apperror.ErrUnauthorized)
			return
		}
		userID, err := tokens.Parse(token)
		if err != nil {
			web.RespondError(w, authResource, apperror.ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// userIDFrom returns the authenticated user ID stored by RequireAuth. The
// second result is false when the request did not pass through RequireAuth.
func userIDFrom(ctx context.Context) (uint64, bool) {
	id, ok := ctx.Value(userIDKey).(uint64)
	return id, ok
}

// Package web provides small HTTP helpers shared by every handler:
// JSON encoding/decoding and a single place to map domain errors to
// HTTP status codes.
package web

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"story-go-mysql/internal/apperror"
)

// errorResponse is the JSON shape returned for every error.
type errorResponse struct {
	Error string `json:"error"`
}

// JSON writes value as a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if value == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("encode response", "error", err)
	}
}

// Error writes a JSON error body with the given status code.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, errorResponse{Error: message})
}

// Decode reads a JSON request body into dst, rejecting unknown fields.
// It returns a ValidationError when the body cannot be decoded.
func Decode(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return apperror.Validation("invalid JSON")
	}
	return nil
}

// RespondError translates a domain error into the appropriate HTTP
// response. resource is used to build a friendly not-found message
// (e.g. "character not found"). Unknown errors are logged and reported
// as 500 without leaking internals to the client.
func RespondError(w http.ResponseWriter, resource string, err error) {
	var validation apperror.ValidationError

	switch {
	case errors.As(err, &validation):
		Error(w, http.StatusBadRequest, validation.Message)
	case errors.Is(err, apperror.ErrNotFound):
		Error(w, http.StatusNotFound, resource+" not found")
	case errors.Is(err, apperror.ErrDuplicateTitle):
		Error(w, http.StatusConflict, "title already exists")
	case errors.Is(err, apperror.ErrInvalidReference):
		Error(w, http.StatusBadRequest, "referenced resource does not exist")
	case errors.Is(err, apperror.ErrUnauthorized):
		Error(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, apperror.ErrDuplicateEmail):
		Error(w, http.StatusConflict, "email already registered")
	default:
		slog.Error("unhandled error", "resource", resource, "error", err)
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}

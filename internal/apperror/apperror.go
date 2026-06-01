// Package apperror defines the domain-level errors used to communicate
// failures between layers without leaking database or HTTP details.
//
// The repository layer translates driver-specific errors (sql.ErrNoRows,
// MySQL duplicate-key, ...) into these values. The service layer adds
// validation errors. The web layer maps them back to HTTP status codes.
package apperror

import "errors"

var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrDuplicateTitle is returned when a unique title constraint is violated.
	ErrDuplicateTitle = errors.New("title already exists")

	// ErrUnauthorized is returned when credentials are missing or invalid.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrDuplicateEmail is returned when an email is already registered.
	ErrDuplicateEmail = errors.New("email already registered")

	// ErrInvalidReference is returned when a request references a related row
	// that does not exist (a foreign-key constraint violation).
	ErrInvalidReference = errors.New("referenced resource does not exist")
)

// ValidationError represents an invalid request that the client can fix.
// It carries a human-readable message that is safe to return to the caller.
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// Validation is a small constructor that keeps call sites concise.
func Validation(message string) ValidationError {
	return ValidationError{Message: message}
}

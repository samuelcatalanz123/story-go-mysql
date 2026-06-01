package repository

import (
	"context"
	"database/sql"
	"time"

	"story-go-mysql/internal/apperror"
)

// RefreshTokenRepository stores refresh tokens. Only the SHA-256 hash of each
// token is persisted, never the token itself: a database leak then can't be
// used to impersonate users.
type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository wires a RefreshTokenRepository to a database handle.
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create stores a refresh token hash for a user with its expiry.
func (r *RefreshTokenRepository) Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)",
		userID, tokenHash, expiresAt,
	)
	return translate(err)
}

// FindValidUser returns the user ID for a token hash that is neither revoked
// nor expired. It returns apperror.ErrUnauthorized when no such token exists.
func (r *RefreshTokenRepository) FindValidUser(ctx context.Context, tokenHash string) (uint64, error) {
	var userID uint64
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id FROM refresh_tokens
		WHERE token_hash = ? AND revoked_at IS NULL AND expires_at > NOW()
	`, tokenHash).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apperror.ErrUnauthorized
		}
		return 0, err
	}
	return userID, nil
}

// Revoke marks a token hash as revoked. Revoking an unknown or already-revoked
// token is a no-op (no error), which keeps logout idempotent.
func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = ? AND revoked_at IS NULL",
		tokenHash,
	)
	return translate(err)
}

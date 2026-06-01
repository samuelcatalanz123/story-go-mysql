package repository

import (
	"context"
	"database/sql"
	"time"

	"story-go-mysql/internal/apperror"
)

// PasswordResetRepository stores one-time password-reset tokens. Only the
// SHA-256 hash is persisted (a DB leak then can't be used to reset passwords).
type PasswordResetRepository struct {
	db *sql.DB
}

// NewPasswordResetRepository wires the repository to a database handle.
func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create stores a reset token hash for a user with its expiry.
func (r *PasswordResetRepository) Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)",
		userID, tokenHash, expiresAt,
	)
	return translate(err)
}

// Consume validates and burns a token in one step: it marks the token used
// only if it is still valid (not used, not expired), then returns the user ID.
// This makes the token single-use even if the link is clicked twice. Returns
// apperror.ErrUnauthorized when the token is invalid, expired or already used.
func (r *PasswordResetRepository) Consume(ctx context.Context, tokenHash string) (uint64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE token_hash = ? AND used_at IS NULL AND expires_at > NOW()
	`, tokenHash)
	if err != nil {
		return 0, translate(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		return 0, apperror.ErrUnauthorized
	}

	var userID uint64
	if err := r.db.QueryRowContext(ctx,
		"SELECT user_id FROM password_reset_tokens WHERE token_hash = ?", tokenHash,
	).Scan(&userID); err != nil {
		return 0, translate(err)
	}
	return userID, nil
}

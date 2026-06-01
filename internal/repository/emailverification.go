package repository

import (
	"context"
	"database/sql"
	"time"

	"story-go-mysql/internal/apperror"
)

// EmailVerificationRepository stores one-time email-verification tokens
// (only their SHA-256 hash), just like password resets.
type EmailVerificationRepository struct {
	db *sql.DB
}

// NewEmailVerificationRepository wires the repository to a database handle.
func NewEmailVerificationRepository(db *sql.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

// Create stores a verification token hash for a user with its expiry.
func (r *EmailVerificationRepository) Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO email_verifications (user_id, token_hash, expires_at) VALUES (?, ?, ?)",
		userID, tokenHash, expiresAt,
	)
	return translate(err)
}

// Consume validates and burns a token atomically (single use), returning the
// user ID. Returns apperror.ErrUnauthorized when invalid, expired or used.
func (r *EmailVerificationRepository) Consume(ctx context.Context, tokenHash string) (uint64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE email_verifications
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
		"SELECT user_id FROM email_verifications WHERE token_hash = ?", tokenHash,
	).Scan(&userID); err != nil {
		return 0, translate(err)
	}
	return userID, nil
}

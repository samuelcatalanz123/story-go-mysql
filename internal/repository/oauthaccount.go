package repository

import (
	"context"
	"database/sql"
)

// OAuthAccountRepository links users to external identity providers (Google,
// etc.) by their immutable provider subject ("sub").
type OAuthAccountRepository struct {
	db *sql.DB
}

// NewOAuthAccountRepository wires an OAuthAccountRepository to a database handle.
func NewOAuthAccountRepository(db *sql.DB) *OAuthAccountRepository {
	return &OAuthAccountRepository{db: db}
}

// FindUserID returns the user ID linked to (provider, subject), or
// apperror.ErrNotFound (via translate) when there is no such link.
func (r *OAuthAccountRepository) FindUserID(ctx context.Context, provider, subject string) (uint64, error) {
	var userID uint64
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id FROM oauth_accounts WHERE provider = ? AND provider_subject = ?
	`, provider, subject).Scan(&userID)
	if err != nil {
		return 0, translate(err)
	}
	return userID, nil
}

// Link records that a user signs in through (provider, subject).
func (r *OAuthAccountRepository) Link(ctx context.Context, userID uint64, provider, subject string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO oauth_accounts (user_id, provider, provider_subject) VALUES (?, ?, ?)",
		userID, provider, subject,
	)
	return translate(err)
}

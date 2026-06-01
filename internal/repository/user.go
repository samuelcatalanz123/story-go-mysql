package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

// UserRepository provides access to the users table.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository wires a UserRepository to a database handle.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a user and returns it. Returns apperror.ErrDuplicateEmail
// when the email is already registered.
func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (model.User, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, passwordHash,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntry {
			return model.User{}, apperror.ErrDuplicateEmail
		}
		return model.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, err
	}
	return r.getByID(ctx, uint64(id))
}

// GetByEmail returns the user and its password hash for login verification.
// Returns apperror.ErrNotFound when no user matches. OAuth-only users have a
// NULL password hash, which COALESCE turns into "" (so password login fails).
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, string, error) {
	var u model.User
	var hash string
	var verifiedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, COALESCE(password_hash, ''), email_verified_at, created_at FROM users WHERE email = ?
	`, email).Scan(&u.ID, &u.Email, &hash, &verifiedAt, &u.CreatedAt)
	if err != nil {
		return model.User{}, "", translate(err)
	}
	if verifiedAt.Valid {
		u.EmailVerifiedAt = &verifiedAt.Time
	}
	return u, hash, nil
}

// CreateOAuthUser creates a user that signs in only through an OAuth provider
// (no password). Returns apperror.ErrDuplicateEmail if the email already exists.
func (r *UserRepository) CreateOAuthUser(ctx context.Context, email string) (model.User, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO users (email, password_hash) VALUES (?, NULL)",
		email,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntry {
			return model.User{}, apperror.ErrDuplicateEmail
		}
		return model.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, err
	}
	return r.getByID(ctx, uint64(id))
}

// GetByID returns the user with the given ID (without the password hash).
func (r *UserRepository) GetByID(ctx context.Context, id uint64) (model.User, error) {
	return r.getByID(ctx, id)
}

// UpdatePassword sets a new password hash for a user, returning
// apperror.ErrNotFound when the user does not exist.
func (r *UserRepository) UpdatePassword(ctx context.Context, id uint64, passwordHash string) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET password_hash = ? WHERE id = ?", passwordHash, id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

func (r *UserRepository) getByID(ctx context.Context, id uint64) (model.User, error) {
	var u model.User
	var verifiedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, email_verified_at, created_at FROM users WHERE id = ?
	`, id).Scan(&u.ID, &u.Email, &verifiedAt, &u.CreatedAt)
	if err != nil {
		return model.User{}, translate(err)
	}
	if verifiedAt.Valid {
		u.EmailVerifiedAt = &verifiedAt.Time
	}
	return u, nil
}

// SetEmailVerified marks the user's email as verified (now), returning
// apperror.ErrNotFound when the user does not exist.
func (r *UserRepository) SetEmailVerified(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET email_verified_at = NOW() WHERE id = ?", id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

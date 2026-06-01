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
// Returns apperror.ErrNotFound when no user matches.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, string, error) {
	var u model.User
	var hash string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at FROM users WHERE email = ?
	`, email).Scan(&u.ID, &u.Email, &hash, &u.CreatedAt)
	if err != nil {
		return model.User{}, "", translate(err)
	}
	return u, hash, nil
}

// GetByID returns the user with the given ID (without the password hash).
func (r *UserRepository) GetByID(ctx context.Context, id uint64) (model.User, error) {
	return r.getByID(ctx, id)
}

func (r *UserRepository) getByID(ctx context.Context, id uint64) (model.User, error) {
	var u model.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, created_at FROM users WHERE id = ?
	`, id).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		return model.User{}, translate(err)
	}
	return u, nil
}

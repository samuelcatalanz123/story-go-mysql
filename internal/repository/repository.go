// Package repository implements the data-access layer. Each repository
// owns the SQL for one aggregate and translates driver-specific errors
// into the domain errors defined in the apperror package.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"

	"story-go-mysql/internal/apperror"
)

// MySQL error numbers we map to domain errors.
const (
	mysqlDuplicateEntry  = 1062 // unique-key violation
	mysqlForeignKeyError = 1452 // foreign-key constraint failure (bad reference)
)

// translate maps low-level driver errors to domain errors. It returns the
// original error unchanged when no mapping applies.
func translate(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return apperror.ErrNotFound
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case mysqlDuplicateEntry:
			return apperror.ErrDuplicateTitle
		case mysqlForeignKeyError:
			return apperror.ErrInvalidReference
		}
	}
	return err
}

// requireAffected converts a "no rows affected" result into
// apperror.ErrNotFound, which callers use to detect missing rows on
// UPDATE/DELETE statements.
func requireAffected(result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// allExist reports whether every id in ids exists in table. An empty slice
// is considered satisfied. The table name is supplied only from trusted
// internal constants, never from user input.
func allExist(ctx context.Context, db *sql.DB, table string, ids []uint64) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	placeholders, args := inPlaceholders(ids)
	query := "SELECT COUNT(DISTINCT id) FROM " + table + " WHERE id IN (" + placeholders + ")"

	var count int
	if err := db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return false, err
	}
	return count == len(ids), nil
}

// inPlaceholders builds the "?,?,?" placeholder string and the matching args
// slice for a SQL IN (...) clause. Callers must guard against an empty slice
// (an empty IN () is invalid SQL).
func inPlaceholders(ids []uint64) (string, []any) {
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	return placeholders, args
}

// buildSearch returns the WHERE clause (with a trailing space) and args for a
// title/text search. An empty query yields an empty clause and nil args, so
// callers can concatenate " ORDER BY ... LIMIT ? OFFSET ?" right after it.
func buildSearch(q string) (string, []any) {
	if q == "" {
		return "", nil
	}
	pattern := "%" + q + "%"
	return "WHERE title LIKE ? OR (text IS NOT NULL AND text LIKE ?) ", []any{pattern, pattern}
}

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

// mysqlDuplicateEntry is the MySQL error number for a unique-key violation.
const mysqlDuplicateEntry = 1062

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
	if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntry {
		return apperror.ErrDuplicateTitle
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

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	query := "SELECT COUNT(DISTINCT id) FROM " + table + " WHERE id IN (" + placeholders + ")"

	var count int
	if err := db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return false, err
	}
	return count == len(ids), nil
}

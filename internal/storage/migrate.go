package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate runs each embedded .sql migration in filename order, exactly once.
// It records every applied migration in the schema_migrations table and skips
// the ones already there. This makes non-idempotent statements (e.g. ALTER
// TABLE ... ADD COLUMN, which has no IF NOT EXISTS in MySQL) safe to ship: a
// migration that already ran is never executed again on a later startup.
func Migrate(ctx context.Context, db *sql.DB) error {
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		applied, err := migrationApplied(ctx, db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		for _, stmt := range splitStatements(string(content)) {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("exec %s: %w", name, err)
			}
		}

		if _, err := db.ExecContext(ctx,
			"INSERT INTO schema_migrations (name) VALUES (?)", name); err != nil {
			return fmt.Errorf("record %s: %w", name, err)
		}
	}
	return nil
}

// ensureMigrationsTable creates the bookkeeping table that tracks which
// migrations have already run. It is itself idempotent (IF NOT EXISTS).
func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	const stmt = `CREATE TABLE IF NOT EXISTS schema_migrations (
		name VARCHAR(255) NOT NULL PRIMARY KEY,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	return nil
}

// migrationApplied reports whether the named migration was already recorded.
func migrationApplied(ctx context.Context, db *sql.DB, name string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name = ?)", name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", name, err)
	}
	return exists, nil
}

// splitStatements divides a SQL file into individual statements by ';',
// discarding empty fragments and surrounding whitespace. (The MySQL driver
// executes one statement per Exec call.)
func splitStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	stmts := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			stmts = append(stmts, s)
		}
	}
	return stmts
}

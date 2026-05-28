package repository

import (
	"context"
	"database/sql"

	"story-go-mysql/internal/model"
)

// LocationRepository provides access to the locations table.
type LocationRepository struct {
	db *sql.DB
}

// NewLocationRepository wires a LocationRepository to a database handle.
func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// Create inserts a location and returns its new ID.
func (r *LocationRepository) Create(ctx context.Context, title string, text *string) (uint64, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO locations (title, text) VALUES (?, ?)",
		title, text,
	)
	if err != nil {
		return 0, translate(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

// GetByID returns a single location or apperror.ErrNotFound.
func (r *LocationRepository) GetByID(ctx context.Context, id uint64) (model.Location, error) {
	var l model.Location
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, text, created_at, updated_at
		FROM locations
		WHERE id = ?
	`, id).Scan(&l.ID, &l.Title, &l.Text, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return model.Location{}, translate(err)
	}
	return l, nil
}

// List returns every location ordered by ID.
func (r *LocationRepository) List(ctx context.Context) ([]model.Location, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, text, created_at, updated_at
		FROM locations
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	locations := []model.Location{}
	for rows.Next() {
		var l model.Location
		if err := rows.Scan(&l.ID, &l.Title, &l.Text, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, rows.Err()
}

// Update modifies a location. It returns apperror.ErrNotFound when no row
// matches the given ID.
func (r *LocationRepository) Update(ctx context.Context, id uint64, title string, text *string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE locations
		SET title = ?, text = ?
		WHERE id = ?
	`, title, text, id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// Delete removes a location, returning apperror.ErrNotFound when missing.
func (r *LocationRepository) Delete(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM locations WHERE id = ?", id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// ExistByIDs reports whether every given location ID exists.
func (r *LocationRepository) ExistByIDs(ctx context.Context, ids []uint64) (bool, error) {
	return allExist(ctx, r.db, "locations", ids)
}

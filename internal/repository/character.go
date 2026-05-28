package repository

import (
	"context"
	"database/sql"

	"story-go-mysql/internal/model"
)

// CharacterRepository provides access to the characters table.
type CharacterRepository struct {
	db *sql.DB
}

// NewCharacterRepository wires a CharacterRepository to a database handle.
func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

// Create inserts a character and returns its new ID.
func (r *CharacterRepository) Create(ctx context.Context, title string, text *string) (uint64, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO characters (title, text) VALUES (?, ?)",
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

// GetByID returns a single character or apperror.ErrNotFound.
func (r *CharacterRepository) GetByID(ctx context.Context, id uint64) (model.Character, error) {
	var c model.Character
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, text, created_at, updated_at
		FROM characters
		WHERE id = ?
	`, id).Scan(&c.ID, &c.Title, &c.Text, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return model.Character{}, translate(err)
	}
	return c, nil
}

// List returns every character ordered by ID.
func (r *CharacterRepository) List(ctx context.Context) ([]model.Character, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, text, created_at, updated_at
		FROM characters
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	characters := []model.Character{}
	for rows.Next() {
		var c model.Character
		if err := rows.Scan(&c.ID, &c.Title, &c.Text, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}
	return characters, rows.Err()
}

// Update modifies a character. It returns apperror.ErrNotFound when no row
// matches the given ID.
func (r *CharacterRepository) Update(ctx context.Context, id uint64, title string, text *string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE characters
		SET title = ?, text = ?
		WHERE id = ?
	`, title, text, id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// Delete removes a character, returning apperror.ErrNotFound when missing.
func (r *CharacterRepository) Delete(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM characters WHERE id = ?", id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// ExistByIDs reports whether every given character ID exists.
func (r *CharacterRepository) ExistByIDs(ctx context.Context, ids []uint64) (bool, error) {
	return allExist(ctx, r.db, "characters", ids)
}

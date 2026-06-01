package repository

import (
	"context"
	"database/sql"

	"story-go-mysql/internal/model"
)

// ConflictRepository provides access to the conflicts table.
type ConflictRepository struct {
	db *sql.DB
}

// NewConflictRepository wires a ConflictRepository to a database handle.
func NewConflictRepository(db *sql.DB) *ConflictRepository {
	return &ConflictRepository{db: db}
}

// Create inserts a conflict and returns its new ID.
func (r *ConflictRepository) Create(ctx context.Context, title string, text *string, sceneID, storyID *uint64) (uint64, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO conflicts (title, text, scene_id, story_id) VALUES (?, ?, ?, ?)",
		title, text, sceneID, storyID,
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

// GetByID returns a single conflict or apperror.ErrNotFound.
func (r *ConflictRepository) GetByID(ctx context.Context, id uint64) (model.Conflict, error) {
	var c model.Conflict
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, text, scene_id, story_id, created_at, updated_at
		FROM conflicts
		WHERE id = ?
	`, id).Scan(&c.ID, &c.Title, &c.Text, &c.SceneID, &c.StoryID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return model.Conflict{}, translate(err)
	}
	return c, nil
}

// List returns a page of conflicts matching q (empty q = no filter),
// ordered by ID, limited to limit rows starting at offset.
func (r *ConflictRepository) List(ctx context.Context, q string, limit, offset int) ([]model.Conflict, error) {
	where, args := buildSearch(q)
	query := "SELECT id, title, text, scene_id, story_id, created_at, updated_at FROM conflicts " +
		where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	conflicts := []model.Conflict{}
	for rows.Next() {
		var c model.Conflict
		if err := rows.Scan(&c.ID, &c.Title, &c.Text, &c.SceneID, &c.StoryID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		conflicts = append(conflicts, c)
	}
	return conflicts, rows.Err()
}

// Count returns the number of conflicts matching q.
func (r *ConflictRepository) Count(ctx context.Context, q string) (int, error) {
	where, args := buildSearch(q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM conflicts "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}

// Update modifies a conflict. It returns apperror.ErrNotFound when no row
// matches the given ID.
func (r *ConflictRepository) Update(ctx context.Context, id uint64, title string, text *string, sceneID, storyID *uint64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE conflicts
		SET title = ?, text = ?, scene_id = ?, story_id = ?
		WHERE id = ?
	`, title, text, sceneID, storyID, id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// Delete removes a conflict, returning apperror.ErrNotFound when missing.
func (r *ConflictRepository) Delete(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM conflicts WHERE id = ?", id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

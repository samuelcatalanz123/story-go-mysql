package repository

import (
	"context"
	"database/sql"

	"story-go-mysql/internal/model"
)

// StoryRepository provides access to the stories table. Every query is
// scoped by user_id so a user can only ever reach their own stories.
type StoryRepository struct {
	db *sql.DB
}

// NewStoryRepository wires a StoryRepository to a database handle.
func NewStoryRepository(db *sql.DB) *StoryRepository {
	return &StoryRepository{db: db}
}

// Create inserts a story owned by userID and returns its new ID.
func (r *StoryRepository) Create(ctx context.Context, userID uint64, title string, text *string) (uint64, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO stories (title, text, user_id) VALUES (?, ?, ?)",
		title, text, userID,
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

// GetByID returns a single story (including its owner) or apperror.ErrNotFound.
// The service compares the owner against the caller before returning it.
func (r *StoryRepository) GetByID(ctx context.Context, id uint64) (model.Story, error) {
	var s model.Story
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, text, user_id, created_at, updated_at
		FROM stories
		WHERE id = ?
	`, id).Scan(&s.ID, &s.Title, &s.Text, &s.UserID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return model.Story{}, translate(err)
	}
	return s, nil
}

// List returns a page of the user's stories matching q (empty q = no filter),
// ordered by ID, limited to limit rows starting at offset.
func (r *StoryRepository) List(ctx context.Context, userID uint64, q string, limit, offset int) ([]model.Story, error) {
	where, args := scopedSearch(userID, q)
	query := "SELECT id, title, text, user_id, created_at, updated_at FROM stories " +
		where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	stories := []model.Story{}
	for rows.Next() {
		var s model.Story
		if err := rows.Scan(&s.ID, &s.Title, &s.Text, &s.UserID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, rows.Err()
}

// Count returns how many of the user's stories match q.
func (r *StoryRepository) Count(ctx context.Context, userID uint64, q string) (int, error) {
	where, args := scopedSearch(userID, q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stories "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}

// Update modifies a story only if it belongs to userID. It returns
// apperror.ErrNotFound when no row matches (missing or owned by someone else).
func (r *StoryRepository) Update(ctx context.Context, id, userID uint64, title string, text *string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE stories
		SET title = ?, text = ?
		WHERE id = ? AND user_id = ?
	`, title, text, id, userID)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// Delete removes a story only if it belongs to userID, returning
// apperror.ErrNotFound when missing or not owned by the caller.
func (r *StoryRepository) Delete(ctx context.Context, id, userID uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM stories WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// scopedSearch builds the WHERE clause (with a trailing space) and args for a
// per-user story query: always filters by user_id, and adds a title/text
// search when q is non-empty.
func scopedSearch(userID uint64, q string) (string, []any) {
	if q == "" {
		return "WHERE user_id = ? ", []any{userID}
	}
	pattern := "%" + q + "%"
	return "WHERE user_id = ? AND (title LIKE ? OR (text IS NOT NULL AND text LIKE ?)) ",
		[]any{userID, pattern, pattern}
}

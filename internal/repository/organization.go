package repository

import (
	"context"
	"database/sql"

	"story-go-mysql/internal/model"
)

// OrganizationRepository provides access to the organizations table.
type OrganizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository wires an OrganizationRepository to a database handle.
func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// Create inserts an organization and returns its new ID.
func (r *OrganizationRepository) Create(ctx context.Context, title string, text *string, storyID *uint64) (uint64, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO organizations (title, text, story_id) VALUES (?, ?, ?)",
		title, text, storyID,
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

// GetByID returns a single organization or apperror.ErrNotFound.
func (r *OrganizationRepository) GetByID(ctx context.Context, id uint64) (model.Organization, error) {
	var o model.Organization
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, text, story_id, created_at, updated_at
		FROM organizations
		WHERE id = ?
	`, id).Scan(&o.ID, &o.Title, &o.Text, &o.StoryID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return model.Organization{}, translate(err)
	}
	return o, nil
}

// List returns a page of organizations matching q (empty q = no filter),
// ordered by ID, limited to limit rows starting at offset.
func (r *OrganizationRepository) List(ctx context.Context, q string, limit, offset int) ([]model.Organization, error) {
	where, args := buildSearch(q)
	query := "SELECT id, title, text, story_id, created_at, updated_at FROM organizations " +
		where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	organizations := []model.Organization{}
	for rows.Next() {
		var o model.Organization
		if err := rows.Scan(&o.ID, &o.Title, &o.Text, &o.StoryID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		organizations = append(organizations, o)
	}
	return organizations, rows.Err()
}

// Count returns the number of organizations matching q.
func (r *OrganizationRepository) Count(ctx context.Context, q string) (int, error) {
	where, args := buildSearch(q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM organizations "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}

// Update modifies an organization. It returns apperror.ErrNotFound when no row
// matches the given ID.
func (r *OrganizationRepository) Update(ctx context.Context, id uint64, title string, text *string, storyID *uint64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE organizations
		SET title = ?, text = ?, story_id = ?
		WHERE id = ?
	`, title, text, storyID, id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// Delete removes an organization, returning apperror.ErrNotFound when missing.
func (r *OrganizationRepository) Delete(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM organizations WHERE id = ?", id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// ExistByIDs reports whether every given organization ID exists.
func (r *OrganizationRepository) ExistByIDs(ctx context.Context, ids []uint64) (bool, error) {
	return allExist(ctx, r.db, "organizations", ids)
}

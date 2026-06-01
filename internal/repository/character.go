package repository

import (
	"context"
	"database/sql"

	"story-go-mysql/internal/model"
)

// CharacterRepository provides access to the characters table and its
// character_organizations join table.
type CharacterRepository struct {
	db *sql.DB
}

// NewCharacterRepository wires a CharacterRepository to a database handle.
func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

// Create inserts a character and its organization links inside a single
// transaction, returning the new character ID.
func (r *CharacterRepository) Create(ctx context.Context, title string, text *string, organizationIDs []uint64) (uint64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // no-op once committed

	result, err := tx.ExecContext(ctx,
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
	characterID := uint64(id)

	if err := insertCharacterOrganizations(ctx, tx, characterID, organizationIDs); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return characterID, nil
}

// GetByID returns a single character, with its organizations, or
// apperror.ErrNotFound.
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

	byCharacter, err := r.organizationsByCharacterIDs(ctx, []uint64{id})
	if err != nil {
		return model.Character{}, err
	}
	c.Organizations = byCharacter[id]
	return c, nil
}

// List returns a page of characters matching q (empty q = no filter),
// ordered by ID, each populated with its organizations. Organizations are
// loaded with a SINGLE extra query for the whole page (not one per row),
// which avoids the N+1 query problem.
func (r *CharacterRepository) List(ctx context.Context, q string, limit, offset int) ([]model.Character, error) {
	where, args := buildSearch(q)
	query := "SELECT id, title, text, created_at, updated_at FROM characters " +
		where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	characters := []model.Character{}
	ids := []uint64{}
	for rows.Next() {
		var c model.Character
		if err := rows.Scan(&c.ID, &c.Title, &c.Text, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		characters = append(characters, c)
		ids = append(ids, c.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// One query for every character's organizations, then attach by ID.
	byCharacter, err := r.organizationsByCharacterIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range characters {
		characters[i].Organizations = byCharacter[characters[i].ID]
	}
	return characters, nil
}

// Count returns the number of characters matching q.
func (r *CharacterRepository) Count(ctx context.Context, q string) (int, error) {
	where, args := buildSearch(q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM characters "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}

// Update modifies a character and replaces its organization links inside a
// transaction. It returns apperror.ErrNotFound when no row matches the ID.
func (r *CharacterRepository) Update(ctx context.Context, id uint64, title string, text *string, organizationIDs []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE characters
		SET title = ?, text = ?
		WHERE id = ?
	`, title, text, id)
	if err != nil {
		return translate(err)
	}
	if err := requireAffected(result); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM character_organizations WHERE character_id = ?", id); err != nil {
		return err
	}
	if err := insertCharacterOrganizations(ctx, tx, id, organizationIDs); err != nil {
		return err
	}
	return tx.Commit()
}

// Delete removes a character, returning apperror.ErrNotFound when missing.
// The character_organizations rows cascade via the foreign key.
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

// organizationsByCharacterIDs returns, for the given character IDs, a map from
// character ID to its organizations. It runs a SINGLE query (WHERE
// character_id IN (...)) instead of one per character — the fix for N+1.
func (r *CharacterRepository) organizationsByCharacterIDs(ctx context.Context, ids []uint64) (map[uint64][]model.Organization, error) {
	result := map[uint64][]model.Organization{}
	if len(ids) == 0 {
		return result, nil
	}

	placeholders, args := inPlaceholders(ids)
	query := `
		SELECT co.character_id, o.id, o.title, o.text, o.story_id, o.created_at, o.updated_at
		FROM organizations o
		JOIN character_organizations co ON o.id = co.organization_id
		WHERE co.character_id IN (` + placeholders + `)
		ORDER BY o.id ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var characterID uint64
		var o model.Organization
		if err := rows.Scan(&characterID, &o.ID, &o.Title, &o.Text, &o.StoryID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		result[characterID] = append(result[characterID], o)
	}
	return result, rows.Err()
}

// insertCharacterOrganizations writes the character_organizations rows for a
// character within the given transaction.
func insertCharacterOrganizations(ctx context.Context, tx *sql.Tx, characterID uint64, organizationIDs []uint64) error {
	for _, orgID := range organizationIDs {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO character_organizations (character_id, organization_id) VALUES (?, ?)",
			characterID, orgID,
		); err != nil {
			return err
		}
	}
	return nil
}

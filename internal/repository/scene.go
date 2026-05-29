package repository

import (
	"context"
	"database/sql"

	"story-go-mysql/internal/model"
)

// SceneData carries the validated fields needed to persist a scene.
// The service layer fills it after validation so the repository only
// deals with concrete values.
type SceneData struct {
	Title         string
	Text          *string
	StartTimeline int
	EndTimeline   int
	CharacterIDs  []uint64
	LocationIDs   []uint64
}

// SceneRepository provides access to the scenes table and its join tables.
type SceneRepository struct {
	db *sql.DB
}

// NewSceneRepository wires a SceneRepository to a database handle.
func NewSceneRepository(db *sql.DB) *SceneRepository {
	return &SceneRepository{db: db}
}

// Create inserts a scene and its character/location links inside a single
// transaction, returning the new scene ID.
func (r *SceneRepository) Create(ctx context.Context, data SceneData) (uint64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // no-op once the transaction is committed

	result, err := tx.ExecContext(ctx, `
		INSERT INTO scenes (title, text, start_timeline, end_timeline)
		VALUES (?, ?, ?, ?)
	`, data.Title, data.Text, data.StartTimeline, data.EndTimeline)
	if err != nil {
		return 0, translate(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	sceneID := uint64(id)

	if err := insertLinks(ctx, tx, sceneID, data); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return sceneID, nil
}

// Update modifies a scene and replaces its links inside a transaction.
// It returns apperror.ErrNotFound when no scene matches the ID.
func (r *SceneRepository) Update(ctx context.Context, id uint64, data SceneData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE scenes
		SET title = ?, text = ?, start_timeline = ?, end_timeline = ?
		WHERE id = ?
	`, data.Title, data.Text, data.StartTimeline, data.EndTimeline, id)
	if err != nil {
		return translate(err)
	}
	if err := requireAffected(result); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM scene_characters WHERE scene_id = ?", id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM scene_locations WHERE scene_id = ?", id); err != nil {
		return err
	}

	if err := insertLinks(ctx, tx, id, data); err != nil {
		return err
	}

	return tx.Commit()
}

// Delete removes a scene. The join tables cascade via foreign keys.
func (r *SceneRepository) Delete(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM scenes WHERE id = ?", id)
	if err != nil {
		return translate(err)
	}
	return requireAffected(result)
}

// GetByID returns a scene with its related characters and locations.
func (r *SceneRepository) GetByID(ctx context.Context, id uint64) (model.Scene, error) {
	var s model.Scene
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, text, start_timeline, end_timeline, created_at, updated_at
		FROM scenes
		WHERE id = ?
	`, id).Scan(&s.ID, &s.Title, &s.Text, &s.StartTimeline, &s.EndTimeline, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return model.Scene{}, translate(err)
	}

	if s.Characters, err = r.characters(ctx, id); err != nil {
		return model.Scene{}, err
	}
	if s.Locations, err = r.locations(ctx, id); err != nil {
		return model.Scene{}, err
	}
	return s, nil
}

// ListIDs returns a page of scene IDs matching q, ordered ascending. The
// service composes full scenes from these IDs via GetByID.
func (r *SceneRepository) ListIDs(ctx context.Context, q string, limit, offset int) ([]uint64, error) {
	where, args := buildSearch(q)
	query := "SELECT id FROM scenes " + where + "ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	ids := []uint64{}
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Count returns the number of scenes matching q.
func (r *SceneRepository) Count(ctx context.Context, q string) (int, error) {
	where, args := buildSearch(q)
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM scenes "+where, args...).Scan(&n); err != nil {
		return 0, translate(err)
	}
	return n, nil
}

func (r *SceneRepository) characters(ctx context.Context, sceneID uint64) ([]model.Character, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.title, c.text, c.created_at, c.updated_at
		FROM characters c
		JOIN scene_characters sc ON c.id = sc.character_id
		WHERE sc.scene_id = ?
		ORDER BY c.id ASC
	`, sceneID)
	if err != nil {
		return nil, err
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

func (r *SceneRepository) locations(ctx context.Context, sceneID uint64) ([]model.Location, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.title, l.text, l.created_at, l.updated_at
		FROM locations l
		JOIN scene_locations sl ON l.id = sl.location_id
		WHERE sl.scene_id = ?
		ORDER BY l.id ASC
	`, sceneID)
	if err != nil {
		return nil, err
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

// insertLinks writes the scene_characters and scene_locations rows for a
// scene within the given transaction.
func insertLinks(ctx context.Context, tx *sql.Tx, sceneID uint64, data SceneData) error {
	for _, characterID := range data.CharacterIDs {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO scene_characters (scene_id, character_id) VALUES (?, ?)",
			sceneID, characterID,
		); err != nil {
			return err
		}
	}
	for _, locationID := range data.LocationIDs {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO scene_locations (scene_id, location_id) VALUES (?, ?)",
			sceneID, locationID,
		); err != nil {
			return err
		}
	}
	return nil
}

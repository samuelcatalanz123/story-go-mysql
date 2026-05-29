package service

import (
	"context"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/repository"
)

// SceneService implements the use cases for scenes. It depends on the
// character and location repositories to validate referenced IDs before
// persisting a scene.
type SceneService struct {
	scenes     *repository.SceneRepository
	characters *repository.CharacterRepository
	locations  *repository.LocationRepository
}

// NewSceneService wires a SceneService to the repositories it needs.
func NewSceneService(
	scenes *repository.SceneRepository,
	characters *repository.CharacterRepository,
	locations *repository.LocationRepository,
) *SceneService {
	return &SceneService{scenes: scenes, characters: characters, locations: locations}
}

// Create validates the request and persists a new scene with its links.
func (s *SceneService) Create(ctx context.Context, req model.SceneRequest) (model.Scene, error) {
	data, err := s.validate(ctx, req)
	if err != nil {
		return model.Scene{}, err
	}
	id, err := s.scenes.Create(ctx, data)
	if err != nil {
		return model.Scene{}, err
	}
	return s.scenes.GetByID(ctx, id)
}

// Update validates the request and replaces the scene and its links.
func (s *SceneService) Update(ctx context.Context, id uint64, req model.SceneRequest) (model.Scene, error) {
	data, err := s.validate(ctx, req)
	if err != nil {
		return model.Scene{}, err
	}
	if err := s.scenes.Update(ctx, id, data); err != nil {
		return model.Scene{}, err
	}
	return s.scenes.GetByID(ctx, id)
}

// Get returns a single scene with its related characters and locations.
func (s *SceneService) Get(ctx context.Context, id uint64) (model.Scene, error) {
	return s.scenes.GetByID(ctx, id)
}

// List returns a page of scenes matching the given params, each populated
// with its relations.
func (s *SceneService) List(ctx context.Context, params model.ListParams) (model.Page[model.Scene], error) {
	p := params.Normalize()
	total, err := s.scenes.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Scene]{}, err
	}
	ids, err := s.scenes.ListIDs(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Scene]{}, err
	}
	scenes := make([]model.Scene, 0, len(ids))
	for _, id := range ids {
		scene, err := s.scenes.GetByID(ctx, id)
		if err != nil {
			return model.Page[model.Scene]{}, err
		}
		scenes = append(scenes, scene)
	}
	return model.Page[model.Scene]{Items: scenes, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
}

// Delete removes a scene by ID.
func (s *SceneService) Delete(ctx context.Context, id uint64) error {
	return s.scenes.Delete(ctx, id)
}

// validate enforces the scene business rules and returns the data ready to
// persist. It checks required fields, the timeline ordering and that every
// referenced character and location exists.
func (s *SceneService) validate(ctx context.Context, req model.SceneRequest) (repository.SceneData, error) {
	if req.Title == "" {
		return repository.SceneData{}, apperror.Validation("title is required")
	}
	if req.StartTimeline == nil {
		return repository.SceneData{}, apperror.Validation("startTimeline is required")
	}
	if req.EndTimeline == nil {
		return repository.SceneData{}, apperror.Validation("endTimeline is required")
	}
	if *req.EndTimeline < *req.StartTimeline {
		return repository.SceneData{}, apperror.Validation("endTimeline must be greater than or equal to startTimeline")
	}

	// Deduplicate so repeated IDs neither cause a false negative in the
	// existence check nor violate the join table's primary key on insert.
	characterIDs := uniqueIDs(req.CharacterIDs)
	locationIDs := uniqueIDs(req.LocationIDs)

	charactersExist, err := s.characters.ExistByIDs(ctx, characterIDs)
	if err != nil {
		return repository.SceneData{}, err
	}
	if !charactersExist {
		return repository.SceneData{}, apperror.Validation("one or more characterIds do not exist")
	}

	locationsExist, err := s.locations.ExistByIDs(ctx, locationIDs)
	if err != nil {
		return repository.SceneData{}, err
	}
	if !locationsExist {
		return repository.SceneData{}, apperror.Validation("one or more locationIds do not exist")
	}

	return repository.SceneData{
		Title:         req.Title,
		Text:          req.Text,
		StartTimeline: *req.StartTimeline,
		EndTimeline:   *req.EndTimeline,
		CharacterIDs:  characterIDs,
		LocationIDs:   locationIDs,
	}, nil
}

// uniqueIDs returns the IDs with duplicates removed, preserving first-seen
// order. A nil or empty input yields an empty slice.
func uniqueIDs(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	result := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

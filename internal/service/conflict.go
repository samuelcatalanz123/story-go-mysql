package service

import (
	"context"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

// conflictStore is the slice of ConflictRepository the service needs.
// Declaring it here keeps the service testable with a fake.
type conflictStore interface {
	Create(ctx context.Context, title string, text *string, sceneID, storyID *uint64) (uint64, error)
	GetByID(ctx context.Context, id uint64) (model.Conflict, error)
	List(ctx context.Context, q string, limit, offset int) ([]model.Conflict, error)
	Count(ctx context.Context, q string) (int, error)
	Update(ctx context.Context, id uint64, title string, text *string, sceneID, storyID *uint64) error
	Delete(ctx context.Context, id uint64) error
}

// ConflictService implements the use cases for conflicts.
type ConflictService struct {
	repo conflictStore
}

// NewConflictService wires a ConflictService to its repository.
func NewConflictService(repo conflictStore) *ConflictService {
	return &ConflictService{repo: repo}
}

// Create validates the request, persists the conflict and returns it.
func (s *ConflictService) Create(ctx context.Context, req model.ConflictRequest) (model.Conflict, error) {
	if req.Title == "" {
		return model.Conflict{}, apperror.Validation("title is required")
	}
	id, err := s.repo.Create(ctx, req.Title, req.Text, req.SceneID, req.StoryID)
	if err != nil {
		return model.Conflict{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// List returns a page of conflicts matching the given params.
func (s *ConflictService) List(ctx context.Context, params model.ListParams) (model.Page[model.Conflict], error) {
	p := params.Normalize()
	total, err := s.repo.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Conflict]{}, err
	}
	items, err := s.repo.List(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Conflict]{}, err
	}
	return model.Page[model.Conflict]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
}

// Get returns a single conflict by ID.
func (s *ConflictService) Get(ctx context.Context, id uint64) (model.Conflict, error) {
	return s.repo.GetByID(ctx, id)
}

// Update validates the request, applies it and returns the updated conflict.
func (s *ConflictService) Update(ctx context.Context, id uint64, req model.ConflictRequest) (model.Conflict, error) {
	if req.Title == "" {
		return model.Conflict{}, apperror.Validation("title is required")
	}
	if err := s.repo.Update(ctx, id, req.Title, req.Text, req.SceneID, req.StoryID); err != nil {
		return model.Conflict{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete removes a conflict by ID.
func (s *ConflictService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

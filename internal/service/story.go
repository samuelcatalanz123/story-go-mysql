package service

import (
	"context"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

// storyStore is the slice of StoryRepository that StoryService needs.
// Declaring it here keeps the service testable with a fake.
type storyStore interface {
	Create(ctx context.Context, userID uint64, title string, text *string) (uint64, error)
	GetByID(ctx context.Context, id uint64) (model.Story, error)
	List(ctx context.Context, userID uint64, q string, limit, offset int) ([]model.Story, error)
	Count(ctx context.Context, userID uint64, q string) (int, error)
	Update(ctx context.Context, id, userID uint64, title string, text *string) error
	Delete(ctx context.Context, id, userID uint64) error
}

// StoryService implements the use cases for stories. Every method takes the
// authenticated userID so a user can only act on their own stories.
type StoryService struct {
	repo storyStore
}

// NewStoryService wires a StoryService to its repository.
func NewStoryService(repo storyStore) *StoryService {
	return &StoryService{repo: repo}
}

// Create validates the request, persists the story for userID and returns it.
func (s *StoryService) Create(ctx context.Context, userID uint64, req model.StoryRequest) (model.Story, error) {
	if req.Title == "" {
		return model.Story{}, apperror.Validation("title is required")
	}
	id, err := s.repo.Create(ctx, userID, req.Title, req.Text)
	if err != nil {
		return model.Story{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// List returns a page of the user's stories matching the given params.
func (s *StoryService) List(ctx context.Context, userID uint64, params model.ListParams) (model.Page[model.Story], error) {
	p := params.Normalize()
	total, err := s.repo.Count(ctx, userID, p.Query)
	if err != nil {
		return model.Page[model.Story]{}, err
	}
	items, err := s.repo.List(ctx, userID, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Story]{}, err
	}
	return model.Page[model.Story]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
}

// Get returns a single story by ID, but only if it belongs to userID.
// Otherwise it returns apperror.ErrNotFound so we never reveal that a story
// owned by another user exists.
func (s *StoryService) Get(ctx context.Context, userID, id uint64) (model.Story, error) {
	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Story{}, err
	}
	if story.UserID != userID {
		return model.Story{}, apperror.ErrNotFound
	}
	return story, nil
}

// Update validates the request and applies it only to a story owned by userID.
func (s *StoryService) Update(ctx context.Context, userID, id uint64, req model.StoryRequest) (model.Story, error) {
	if req.Title == "" {
		return model.Story{}, apperror.Validation("title is required")
	}
	if err := s.repo.Update(ctx, id, userID, req.Title, req.Text); err != nil {
		return model.Story{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete removes a story owned by userID.
func (s *StoryService) Delete(ctx context.Context, userID, id uint64) error {
	return s.repo.Delete(ctx, id, userID)
}

// Package service holds the business logic. It validates input, applies
// rules and orchestrates repositories, returning domain types and the
// errors defined in the apperror package.
package service

import (
	"context"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/repository"
)

// CharacterService implements the use cases for characters.
type CharacterService struct {
	repo *repository.CharacterRepository
}

// NewCharacterService wires a CharacterService to its repository.
func NewCharacterService(repo *repository.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

// Create validates the request, persists the character and returns it.
func (s *CharacterService) Create(ctx context.Context, req model.CharacterRequest) (model.Character, error) {
	if req.Title == "" {
		return model.Character{}, apperror.Validation("title is required")
	}
	id, err := s.repo.Create(ctx, req.Title, req.Text)
	if err != nil {
		return model.Character{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// List returns every character.
func (s *CharacterService) List(ctx context.Context) ([]model.Character, error) {
	return s.repo.List(ctx)
}

// Get returns a single character by ID.
func (s *CharacterService) Get(ctx context.Context, id uint64) (model.Character, error) {
	return s.repo.GetByID(ctx, id)
}

// Update validates the request, applies it and returns the updated character.
func (s *CharacterService) Update(ctx context.Context, id uint64, req model.CharacterRequest) (model.Character, error) {
	if req.Title == "" {
		return model.Character{}, apperror.Validation("title is required")
	}
	if err := s.repo.Update(ctx, id, req.Title, req.Text); err != nil {
		return model.Character{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete removes a character by ID.
func (s *CharacterService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

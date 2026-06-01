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

// CharacterService implements the use cases for characters. It depends on the
// organization repository to validate referenced organization IDs.
type CharacterService struct {
	repo          *repository.CharacterRepository
	organizations *repository.OrganizationRepository
}

// NewCharacterService wires a CharacterService to the repositories it needs.
func NewCharacterService(repo *repository.CharacterRepository, organizations *repository.OrganizationRepository) *CharacterService {
	return &CharacterService{repo: repo, organizations: organizations}
}

// Create validates the request, persists the character and returns it.
func (s *CharacterService) Create(ctx context.Context, req model.CharacterRequest) (model.Character, error) {
	if req.Title == "" {
		return model.Character{}, apperror.Validation("title is required")
	}
	organizationIDs, err := s.validOrganizationIDs(ctx, req.OrganizationIDs)
	if err != nil {
		return model.Character{}, err
	}
	id, err := s.repo.Create(ctx, req.Title, req.Text, organizationIDs)
	if err != nil {
		return model.Character{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// List returns a page of characters matching the given params.
func (s *CharacterService) List(ctx context.Context, params model.ListParams) (model.Page[model.Character], error) {
	p := params.Normalize()
	total, err := s.repo.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Character]{}, err
	}
	items, err := s.repo.List(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Character]{}, err
	}
	return model.Page[model.Character]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
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
	organizationIDs, err := s.validOrganizationIDs(ctx, req.OrganizationIDs)
	if err != nil {
		return model.Character{}, err
	}
	if err := s.repo.Update(ctx, id, req.Title, req.Text, organizationIDs); err != nil {
		return model.Character{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete removes a character by ID.
func (s *CharacterService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// validOrganizationIDs deduplicates the requested organization IDs and checks
// that every one exists, returning a ValidationError otherwise.
func (s *CharacterService) validOrganizationIDs(ctx context.Context, ids []uint64) ([]uint64, error) {
	unique := uniqueIDs(ids)
	exist, err := s.organizations.ExistByIDs(ctx, unique)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, apperror.Validation("one or more organizationIds do not exist")
	}
	return unique, nil
}

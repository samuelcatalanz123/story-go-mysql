package service

import (
	"context"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

// organizationStore is the slice of OrganizationRepository that the service
// needs. Declaring it here keeps the service testable with a fake.
type organizationStore interface {
	Create(ctx context.Context, title string, text *string, storyID *uint64) (uint64, error)
	GetByID(ctx context.Context, id uint64) (model.Organization, error)
	List(ctx context.Context, q string, limit, offset int) ([]model.Organization, error)
	Count(ctx context.Context, q string) (int, error)
	Update(ctx context.Context, id uint64, title string, text *string, storyID *uint64) error
	Delete(ctx context.Context, id uint64) error
}

// OrganizationService implements the use cases for organizations.
type OrganizationService struct {
	repo organizationStore
}

// NewOrganizationService wires an OrganizationService to its repository.
func NewOrganizationService(repo organizationStore) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// Create validates the request, persists the organization and returns it.
func (s *OrganizationService) Create(ctx context.Context, req model.OrganizationRequest) (model.Organization, error) {
	if req.Title == "" {
		return model.Organization{}, apperror.Validation("title is required")
	}
	id, err := s.repo.Create(ctx, req.Title, req.Text, req.StoryID)
	if err != nil {
		return model.Organization{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// List returns a page of organizations matching the given params.
func (s *OrganizationService) List(ctx context.Context, params model.ListParams) (model.Page[model.Organization], error) {
	p := params.Normalize()
	total, err := s.repo.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Organization]{}, err
	}
	items, err := s.repo.List(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Organization]{}, err
	}
	return model.Page[model.Organization]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize}, nil
}

// Get returns a single organization by ID.
func (s *OrganizationService) Get(ctx context.Context, id uint64) (model.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

// Update validates the request, applies it and returns the updated organization.
func (s *OrganizationService) Update(ctx context.Context, id uint64, req model.OrganizationRequest) (model.Organization, error) {
	if req.Title == "" {
		return model.Organization{}, apperror.Validation("title is required")
	}
	if err := s.repo.Update(ctx, id, req.Title, req.Text, req.StoryID); err != nil {
		return model.Organization{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete removes an organization by ID.
func (s *OrganizationService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

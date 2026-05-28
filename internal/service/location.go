package service

import (
	"context"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/repository"
)

// LocationService implements the use cases for locations.
type LocationService struct {
	repo *repository.LocationRepository
}

// NewLocationService wires a LocationService to its repository.
func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

// Create validates the request, persists the location and returns it.
func (s *LocationService) Create(ctx context.Context, req model.LocationRequest) (model.Location, error) {
	if req.Title == "" {
		return model.Location{}, apperror.Validation("title is required")
	}
	id, err := s.repo.Create(ctx, req.Title, req.Text)
	if err != nil {
		return model.Location{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// List returns every location.
func (s *LocationService) List(ctx context.Context) ([]model.Location, error) {
	return s.repo.List(ctx)
}

// Get returns a single location by ID.
func (s *LocationService) Get(ctx context.Context, id uint64) (model.Location, error) {
	return s.repo.GetByID(ctx, id)
}

// Update validates the request, applies it and returns the updated location.
func (s *LocationService) Update(ctx context.Context, id uint64, req model.LocationRequest) (model.Location, error) {
	if req.Title == "" {
		return model.Location{}, apperror.Validation("title is required")
	}
	if err := s.repo.Update(ctx, id, req.Title, req.Text); err != nil {
		return model.Location{}, err
	}
	return s.repo.GetByID(ctx, id)
}

// Delete removes a location by ID.
func (s *LocationService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

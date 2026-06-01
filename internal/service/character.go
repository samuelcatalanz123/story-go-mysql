// Package service holds the business logic. It validates input, applies
// rules and orchestrates repositories, returning domain types and the
// errors defined in the apperror package.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/cache"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/repository"
)

// characterListPrefix namespaces the cached character-list entries so we can
// invalidate them all at once after a write.
const characterListPrefix = "characters:list:"

// characterListTTL is how long a cached list stays fresh.
const characterListTTL = 60 * time.Second

// CharacterService implements the use cases for characters. It depends on the
// organization repository to validate referenced organization IDs, and on a
// cache to speed up list reads (cache-aside).
type CharacterService struct {
	repo          *repository.CharacterRepository
	organizations *repository.OrganizationRepository
	cache         cache.Cache
}

// NewCharacterService wires a CharacterService to the repositories it needs.
// Pass cache.Noop{} to disable caching.
func NewCharacterService(repo *repository.CharacterRepository, organizations *repository.OrganizationRepository, c cache.Cache) *CharacterService {
	return &CharacterService{repo: repo, organizations: organizations, cache: c}
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
	s.invalidateLists(ctx)
	return s.repo.GetByID(ctx, id)
}

// List returns a page of characters matching the given params, using the
// cache-aside pattern: try the cache first, and on a miss read from MySQL and
// store the result with a TTL.
func (s *CharacterService) List(ctx context.Context, params model.ListParams) (model.Page[model.Character], error) {
	p := params.Normalize()
	key := fmt.Sprintf("%s%d:%d:%s", characterListPrefix, p.Page, p.PageSize, p.Query)

	// 1) ¿Está en caché?
	if data, ok, err := s.cache.Get(ctx, key); err == nil && ok {
		var page model.Page[model.Character]
		if json.Unmarshal(data, &page) == nil {
			return page, nil
		}
	}

	// 2) Miss: leer de MySQL.
	total, err := s.repo.Count(ctx, p.Query)
	if err != nil {
		return model.Page[model.Character]{}, err
	}
	items, err := s.repo.List(ctx, p.Query, p.Limit(), p.Offset())
	if err != nil {
		return model.Page[model.Character]{}, err
	}
	page := model.Page[model.Character]{Items: items, Total: total, Page: p.Page, PageSize: p.PageSize}

	// 3) Guardar en caché con expiración (si falla, no pasa nada: ya tenemos el dato).
	if data, err := json.Marshal(page); err == nil {
		if err := s.cache.Set(ctx, key, data, characterListTTL); err != nil {
			slog.Warn("no se pudo cachear la lista de personajes", "error", err)
		}
	}
	return page, nil
}

// invalidateLists borra todas las listas de personajes cacheadas. Se llama tras
// cualquier escritura para que nadie lea datos obsoletos.
func (s *CharacterService) invalidateLists(ctx context.Context) {
	if err := s.cache.DelByPrefix(ctx, characterListPrefix); err != nil {
		slog.Warn("no se pudo invalidar la caché de personajes", "error", err)
	}
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
	s.invalidateLists(ctx)
	return s.repo.GetByID(ctx, id)
}

// Delete removes a character by ID.
func (s *CharacterService) Delete(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateLists(ctx)
	return nil
}

// SetAvatar stores the avatar path and returns the updated character.
func (s *CharacterService) SetAvatar(ctx context.Context, id uint64, path string) (model.Character, error) {
	if err := s.repo.SetAvatar(ctx, id, path); err != nil {
		return model.Character{}, err
	}
	s.invalidateLists(ctx)
	return s.repo.GetByID(ctx, id)
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

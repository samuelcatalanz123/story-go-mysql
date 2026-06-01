package service

import (
	"context"
	"errors"
	"testing"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

// fakeConflictStore is an in-memory conflictStore for testing, without a DB.
type fakeConflictStore struct {
	conflicts map[uint64]model.Conflict
	nextID    uint64
}

func newFakeConflictStore() *fakeConflictStore {
	return &fakeConflictStore{conflicts: map[uint64]model.Conflict{}}
}

func (f *fakeConflictStore) Create(_ context.Context, title string, text *string, sceneID, storyID *uint64) (uint64, error) {
	f.nextID++
	f.conflicts[f.nextID] = model.Conflict{ID: f.nextID, Title: title, Text: text, SceneID: sceneID, StoryID: storyID}
	return f.nextID, nil
}

func (f *fakeConflictStore) GetByID(_ context.Context, id uint64) (model.Conflict, error) {
	c, ok := f.conflicts[id]
	if !ok {
		return model.Conflict{}, apperror.ErrNotFound
	}
	return c, nil
}

func (f *fakeConflictStore) List(_ context.Context, _ string, _, _ int) ([]model.Conflict, error) {
	out := []model.Conflict{}
	for _, c := range f.conflicts {
		out = append(out, c)
	}
	return out, nil
}

func (f *fakeConflictStore) Count(_ context.Context, _ string) (int, error) {
	return len(f.conflicts), nil
}

func (f *fakeConflictStore) Update(_ context.Context, id uint64, title string, text *string, sceneID, storyID *uint64) error {
	c, ok := f.conflicts[id]
	if !ok {
		return apperror.ErrNotFound
	}
	c.Title, c.Text, c.SceneID, c.StoryID = title, text, sceneID, storyID
	f.conflicts[id] = c
	return nil
}

func (f *fakeConflictStore) Delete(_ context.Context, id uint64) error {
	if _, ok := f.conflicts[id]; !ok {
		return apperror.ErrNotFound
	}
	delete(f.conflicts, id)
	return nil
}

func TestConflictCreateRequiresTitle(t *testing.T) {
	svc := NewConflictService(newFakeConflictStore())
	_, err := svc.Create(context.Background(), model.ConflictRequest{Title: ""})
	var v apperror.ValidationError
	if !errors.As(err, &v) {
		t.Fatalf("esperaba ValidationError, obtuve %v", err)
	}
}

func TestConflictCreateAndGet(t *testing.T) {
	svc := NewConflictService(newFakeConflictStore())
	sceneID := uint64(5)
	created, err := svc.Create(context.Background(), model.ConflictRequest{Title: "La Batalla Final", SceneID: &sceneID})
	if err != nil {
		t.Fatal(err)
	}
	got, err := svc.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "La Batalla Final" {
		t.Fatalf("título inesperado: %s", got.Title)
	}
	if got.SceneID == nil || *got.SceneID != 5 {
		t.Fatalf("esperaba sceneID=5, obtuve %v", got.SceneID)
	}
}

func TestConflictUpdateMissing(t *testing.T) {
	svc := NewConflictService(newFakeConflictStore())
	_, err := svc.Update(context.Background(), 999, model.ConflictRequest{Title: "x"})
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, obtuve %v", err)
	}
}

func TestConflictDelete(t *testing.T) {
	svc := NewConflictService(newFakeConflictStore())
	created, _ := svc.Create(context.Background(), model.ConflictRequest{Title: "Disputa"})
	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Get(context.Background(), created.ID); !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("esperaba que ya no existiera, obtuve %v", err)
	}
}

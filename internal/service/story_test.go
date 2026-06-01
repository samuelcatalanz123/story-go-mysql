package service

import (
	"context"
	"errors"
	"testing"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

// fakeStoryStore is an in-memory storyStore for testing, without a database.
type fakeStoryStore struct {
	stories map[uint64]model.Story
	nextID  uint64
}

func newFakeStoryStore() *fakeStoryStore {
	return &fakeStoryStore{stories: map[uint64]model.Story{}}
}

func (f *fakeStoryStore) Create(_ context.Context, userID uint64, title string, text *string) (uint64, error) {
	f.nextID++
	f.stories[f.nextID] = model.Story{ID: f.nextID, Title: title, Text: text, UserID: userID}
	return f.nextID, nil
}

func (f *fakeStoryStore) GetByID(_ context.Context, id uint64) (model.Story, error) {
	s, ok := f.stories[id]
	if !ok {
		return model.Story{}, apperror.ErrNotFound
	}
	return s, nil
}

func (f *fakeStoryStore) List(_ context.Context, userID uint64, _ string, _, _ int) ([]model.Story, error) {
	out := []model.Story{}
	for _, s := range f.stories {
		if s.UserID == userID {
			out = append(out, s)
		}
	}
	return out, nil
}

func (f *fakeStoryStore) Count(_ context.Context, userID uint64, _ string) (int, error) {
	n := 0
	for _, s := range f.stories {
		if s.UserID == userID {
			n++
		}
	}
	return n, nil
}

func (f *fakeStoryStore) Update(_ context.Context, id, userID uint64, title string, text *string) error {
	s, ok := f.stories[id]
	if !ok || s.UserID != userID {
		return apperror.ErrNotFound
	}
	s.Title, s.Text = title, text
	f.stories[id] = s
	return nil
}

func (f *fakeStoryStore) Delete(_ context.Context, id, userID uint64) error {
	s, ok := f.stories[id]
	if !ok || s.UserID != userID {
		return apperror.ErrNotFound
	}
	delete(f.stories, id)
	return nil
}

func TestStoryCreateRequiresTitle(t *testing.T) {
	svc := NewStoryService(newFakeStoryStore())
	_, err := svc.Create(context.Background(), 1, model.StoryRequest{Title: ""})
	var v apperror.ValidationError
	if !errors.As(err, &v) {
		t.Fatalf("esperaba ValidationError, obtuve %v", err)
	}
}

func TestStoryCreateAssignsOwner(t *testing.T) {
	svc := NewStoryService(newFakeStoryStore())
	story, err := svc.Create(context.Background(), 42, model.StoryRequest{Title: "Mi historia"})
	if err != nil {
		t.Fatal(err)
	}
	if story.UserID != 42 {
		t.Fatalf("esperaba dueño 42, obtuve %d", story.UserID)
	}
	if story.Title != "Mi historia" {
		t.Fatalf("título inesperado: %s", story.Title)
	}
}

func TestStoryGetHidesOtherUsersStory(t *testing.T) {
	svc := NewStoryService(newFakeStoryStore())
	// El usuario 1 crea una historia.
	created, _ := svc.Create(context.Background(), 1, model.StoryRequest{Title: "Privada"})
	// El usuario 2 intenta verla: debe parecer "no encontrada".
	_, err := svc.Get(context.Background(), 2, created.ID)
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound para otro usuario, obtuve %v", err)
	}
	// El dueño sí la ve.
	if _, err := svc.Get(context.Background(), 1, created.ID); err != nil {
		t.Fatalf("el dueño debería verla, obtuve %v", err)
	}
}

func TestStoryDeleteOnlyByOwner(t *testing.T) {
	svc := NewStoryService(newFakeStoryStore())
	created, _ := svc.Create(context.Background(), 1, model.StoryRequest{Title: "Mía"})
	if err := svc.Delete(context.Background(), 2, created.ID); !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("otro usuario no debería borrarla, obtuve %v", err)
	}
	if err := svc.Delete(context.Background(), 1, created.ID); err != nil {
		t.Fatalf("el dueño debería borrarla, obtuve %v", err)
	}
}

func TestStoryListOnlyOwnStories(t *testing.T) {
	svc := NewStoryService(newFakeStoryStore())
	_, _ = svc.Create(context.Background(), 1, model.StoryRequest{Title: "A"})
	_, _ = svc.Create(context.Background(), 1, model.StoryRequest{Title: "B"})
	_, _ = svc.Create(context.Background(), 2, model.StoryRequest{Title: "C"})

	page, err := svc.List(context.Background(), 1, model.ListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 2 {
		t.Fatalf("esperaba 2 historias del usuario 1, obtuve %d", page.Total)
	}
}

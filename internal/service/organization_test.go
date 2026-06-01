package service

import (
	"context"
	"errors"
	"testing"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
)

// fakeOrgStore is an in-memory organizationStore for testing, without a DB.
type fakeOrgStore struct {
	orgs   map[uint64]model.Organization
	nextID uint64
}

func newFakeOrgStore() *fakeOrgStore {
	return &fakeOrgStore{orgs: map[uint64]model.Organization{}}
}

func (f *fakeOrgStore) Create(_ context.Context, title string, text *string, storyID *uint64) (uint64, error) {
	f.nextID++
	f.orgs[f.nextID] = model.Organization{ID: f.nextID, Title: title, Text: text, StoryID: storyID}
	return f.nextID, nil
}

func (f *fakeOrgStore) GetByID(_ context.Context, id uint64) (model.Organization, error) {
	o, ok := f.orgs[id]
	if !ok {
		return model.Organization{}, apperror.ErrNotFound
	}
	return o, nil
}

func (f *fakeOrgStore) List(_ context.Context, _ string, _, _ int) ([]model.Organization, error) {
	out := []model.Organization{}
	for _, o := range f.orgs {
		out = append(out, o)
	}
	return out, nil
}

func (f *fakeOrgStore) Count(_ context.Context, _ string) (int, error) {
	return len(f.orgs), nil
}

func (f *fakeOrgStore) Update(_ context.Context, id uint64, title string, text *string, storyID *uint64) error {
	o, ok := f.orgs[id]
	if !ok {
		return apperror.ErrNotFound
	}
	o.Title, o.Text, o.StoryID = title, text, storyID
	f.orgs[id] = o
	return nil
}

func (f *fakeOrgStore) Delete(_ context.Context, id uint64) error {
	if _, ok := f.orgs[id]; !ok {
		return apperror.ErrNotFound
	}
	delete(f.orgs, id)
	return nil
}

func TestOrganizationCreateRequiresTitle(t *testing.T) {
	svc := NewOrganizationService(newFakeOrgStore())
	_, err := svc.Create(context.Background(), model.OrganizationRequest{Title: ""})
	var v apperror.ValidationError
	if !errors.As(err, &v) {
		t.Fatalf("esperaba ValidationError, obtuve %v", err)
	}
}

func TestOrganizationCreateAndGet(t *testing.T) {
	svc := NewOrganizationService(newFakeOrgStore())
	created, err := svc.Create(context.Background(), model.OrganizationRequest{Title: "La Tripulación"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := svc.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "La Tripulación" {
		t.Fatalf("título inesperado: %s", got.Title)
	}
}

func TestOrganizationGetMissing(t *testing.T) {
	svc := NewOrganizationService(newFakeOrgStore())
	_, err := svc.Get(context.Background(), 999)
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("esperaba ErrNotFound, obtuve %v", err)
	}
}

func TestOrganizationDelete(t *testing.T) {
	svc := NewOrganizationService(newFakeOrgStore())
	created, _ := svc.Create(context.Background(), model.OrganizationRequest{Title: "Gremio"})
	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Get(context.Background(), created.ID); !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("esperaba que ya no existiera, obtuve %v", err)
	}
}

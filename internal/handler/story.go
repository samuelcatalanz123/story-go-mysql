package handler

import (
	"net/http"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const storyResource = "story"

// StoryHandler exposes the story endpoints. All of them run behind
// RequireAuth, so the authenticated user ID is always present in the context.
type StoryHandler struct {
	svc *service.StoryService
}

// NewStoryHandler wires a StoryHandler to its service.
func NewStoryHandler(svc *service.StoryService) *StoryHandler {
	return &StoryHandler{svc: svc}
}

func (h *StoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFrom(r.Context())
	if !ok {
		web.RespondError(w, storyResource, apperror.ErrUnauthorized)
		return
	}
	var req model.StoryRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, storyResource, err)
		return
	}
	story, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		web.RespondError(w, storyResource, err)
		return
	}
	web.JSON(w, http.StatusCreated, story)
}

func (h *StoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFrom(r.Context())
	if !ok {
		web.RespondError(w, storyResource, apperror.ErrUnauthorized)
		return
	}
	page, err := h.svc.List(r.Context(), userID, parseListParams(r))
	if err != nil {
		web.RespondError(w, storyResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}

func (h *StoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFrom(r.Context())
	if !ok {
		web.RespondError(w, storyResource, apperror.ErrUnauthorized)
		return
	}
	id, ok := parseID(w, r, storyResource)
	if !ok {
		return
	}
	story, err := h.svc.Get(r.Context(), userID, id)
	if err != nil {
		web.RespondError(w, storyResource, err)
		return
	}
	web.JSON(w, http.StatusOK, story)
}

func (h *StoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFrom(r.Context())
	if !ok {
		web.RespondError(w, storyResource, apperror.ErrUnauthorized)
		return
	}
	id, ok := parseID(w, r, storyResource)
	if !ok {
		return
	}
	var req model.StoryRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, storyResource, err)
		return
	}
	story, err := h.svc.Update(r.Context(), userID, id, req)
	if err != nil {
		web.RespondError(w, storyResource, err)
		return
	}
	web.JSON(w, http.StatusOK, story)
}

func (h *StoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFrom(r.Context())
	if !ok {
		web.RespondError(w, storyResource, apperror.ErrUnauthorized)
		return
	}
	id, ok := parseID(w, r, storyResource)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), userID, id); err != nil {
		web.RespondError(w, storyResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "story deleted"})
}

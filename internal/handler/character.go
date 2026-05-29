package handler

import (
	"net/http"

	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const characterResource = "character"

// CharacterHandler exposes the character endpoints.
type CharacterHandler struct {
	svc *service.CharacterService
}

// NewCharacterHandler wires a CharacterHandler to its service.
func NewCharacterHandler(svc *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{svc: svc}
}

func (h *CharacterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CharacterRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	character, err := h.svc.Create(r.Context(), req)
	if err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	web.JSON(w, http.StatusCreated, character)
}

func (h *CharacterHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := h.svc.List(r.Context(), parseListParams(r))
	if err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}

func (h *CharacterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, characterResource)
	if !ok {
		return
	}
	character, err := h.svc.Get(r.Context(), id)
	if err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	web.JSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, characterResource)
	if !ok {
		return
	}
	var req model.CharacterRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	character, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	web.JSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, characterResource)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		web.RespondError(w, characterResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "character deleted"})
}

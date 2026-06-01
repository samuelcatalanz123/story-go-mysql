package handler

import (
	"net/http"

	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const conflictResource = "conflict"

// ConflictHandler exposes the conflict endpoints.
type ConflictHandler struct {
	svc *service.ConflictService
}

// NewConflictHandler wires a ConflictHandler to its service.
func NewConflictHandler(svc *service.ConflictService) *ConflictHandler {
	return &ConflictHandler{svc: svc}
}

func (h *ConflictHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.ConflictRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, conflictResource, err)
		return
	}
	conflict, err := h.svc.Create(r.Context(), req)
	if err != nil {
		web.RespondError(w, conflictResource, err)
		return
	}
	web.JSON(w, http.StatusCreated, conflict)
}

func (h *ConflictHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := h.svc.List(r.Context(), parseListParams(r))
	if err != nil {
		web.RespondError(w, conflictResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}

func (h *ConflictHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, conflictResource)
	if !ok {
		return
	}
	conflict, err := h.svc.Get(r.Context(), id)
	if err != nil {
		web.RespondError(w, conflictResource, err)
		return
	}
	web.JSON(w, http.StatusOK, conflict)
}

func (h *ConflictHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, conflictResource)
	if !ok {
		return
	}
	var req model.ConflictRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, conflictResource, err)
		return
	}
	conflict, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		web.RespondError(w, conflictResource, err)
		return
	}
	web.JSON(w, http.StatusOK, conflict)
}

func (h *ConflictHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, conflictResource)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		web.RespondError(w, conflictResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "conflict deleted"})
}

package handler

import (
	"net/http"

	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const locationResource = "location"

// LocationHandler exposes the location endpoints.
type LocationHandler struct {
	svc *service.LocationService
}

// NewLocationHandler wires a LocationHandler to its service.
func NewLocationHandler(svc *service.LocationService) *LocationHandler {
	return &LocationHandler{svc: svc}
}

func (h *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.LocationRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	location, err := h.svc.Create(r.Context(), req)
	if err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	web.JSON(w, http.StatusCreated, location)
}

func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	locations, err := h.svc.List(r.Context())
	if err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, locations)
}

func (h *LocationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, locationResource)
	if !ok {
		return
	}
	location, err := h.svc.Get(r.Context(), id)
	if err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, location)
}

func (h *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, locationResource)
	if !ok {
		return
	}
	var req model.LocationRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	location, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, location)
}

func (h *LocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, locationResource)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		web.RespondError(w, locationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "location deleted"})
}

package handler

import (
	"net/http"

	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const sceneResource = "scene"

// SceneHandler exposes the scene endpoints.
type SceneHandler struct {
	svc *service.SceneService
}

// NewSceneHandler wires a SceneHandler to its service.
func NewSceneHandler(svc *service.SceneService) *SceneHandler {
	return &SceneHandler{svc: svc}
}

func (h *SceneHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.SceneRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	scene, err := h.svc.Create(r.Context(), req)
	if err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	web.JSON(w, http.StatusCreated, scene)
}

func (h *SceneHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := h.svc.List(r.Context(), parseListParams(r))
	if err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}

func (h *SceneHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, sceneResource)
	if !ok {
		return
	}
	scene, err := h.svc.Get(r.Context(), id)
	if err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	web.JSON(w, http.StatusOK, scene)
}

func (h *SceneHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, sceneResource)
	if !ok {
		return
	}
	var req model.SceneRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	scene, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	web.JSON(w, http.StatusOK, scene)
}

func (h *SceneHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, sceneResource)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		web.RespondError(w, sceneResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "scene deleted"})
}

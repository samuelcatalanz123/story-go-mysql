package handler

import (
	"net/http"

	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const organizationResource = "organization"

// OrganizationHandler exposes the organization endpoints.
type OrganizationHandler struct {
	svc *service.OrganizationService
}

// NewOrganizationHandler wires an OrganizationHandler to its service.
func NewOrganizationHandler(svc *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{svc: svc}
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.OrganizationRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, organizationResource, err)
		return
	}
	organization, err := h.svc.Create(r.Context(), req)
	if err != nil {
		web.RespondError(w, organizationResource, err)
		return
	}
	web.JSON(w, http.StatusCreated, organization)
}

func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := h.svc.List(r.Context(), parseListParams(r))
	if err != nil {
		web.RespondError(w, organizationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, page)
}

func (h *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, organizationResource)
	if !ok {
		return
	}
	organization, err := h.svc.Get(r.Context(), id)
	if err != nil {
		web.RespondError(w, organizationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, organization)
}

func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, organizationResource)
	if !ok {
		return
	}
	var req model.OrganizationRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, organizationResource, err)
		return
	}
	organization, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		web.RespondError(w, organizationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, organization)
}

func (h *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, organizationResource)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		web.RespondError(w, organizationResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "organization deleted"})
}

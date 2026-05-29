package handler

import (
	"net/http"

	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const authResource = "auth"

// AuthHandler exposes the registration and login endpoints.
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler wires an AuthHandler to its service.
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	res, err := h.svc.Register(r.Context(), req)
	if err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	web.JSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	res, err := h.svc.Login(r.Context(), req)
	if err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	web.JSON(w, http.StatusOK, res)
}

package handler

import (
	"net/http"

	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

// PasswordHandler exposes the forgot/reset password endpoints.
type PasswordHandler struct {
	svc *service.PasswordResetService
}

// NewPasswordHandler wires a PasswordHandler to its service.
func NewPasswordHandler(svc *service.PasswordResetService) *PasswordHandler {
	return &PasswordHandler{svc: svc}
}

// Forgot handles POST /auth/forgot-password. It always responds 200 with the
// same vague message so an attacker can't tell which emails are registered.
func (h *PasswordHandler) Forgot(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	if err := h.svc.ForgotPassword(r.Context(), req.Email); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{
		"message": "Si existe una cuenta con ese email, te enviamos un enlace para restablecer la contraseña.",
	})
}

// Reset handles POST /auth/reset-password.
func (h *PasswordHandler) Reset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	if err := h.svc.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "Contraseña actualizada"})
}

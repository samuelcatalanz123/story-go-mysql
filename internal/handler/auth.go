package handler

import (
	"log/slog"
	"net/http"
	"time"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/web"
)

const authResource = "auth"

// refreshCookieName / refreshCookiePath: the refresh token lives in an
// HttpOnly cookie scoped to the auth routes, so it is sent only to
// /api/auth/* and is invisible to JavaScript (immune to XSS token theft).
const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/auth"
)

// AuthHandler exposes registration, login, refresh, logout, Google OAuth and
// email verification.
type AuthHandler struct {
	svc          *service.AuthService
	oauth        *service.OAuthService
	verification *service.EmailVerificationService
	refreshTTL   time.Duration
}

// NewAuthHandler wires an AuthHandler to its services and the refresh-token
// lifetime (used for the cookie's Max-Age).
func NewAuthHandler(svc *service.AuthService, oauth *service.OAuthService, verification *service.EmailVerificationService, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{svc: svc, oauth: oauth, verification: verification, refreshTTL: refreshTTL}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	res, refresh, err := h.svc.Register(r.Context(), req)
	if err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	// Envía el correo de verificación (best-effort: si el email falla, el
	// registro igual fue exitoso; el usuario puede reenviarlo después).
	if err := h.verification.Send(r.Context(), res.User.ID); err != nil {
		slog.Warn("no se pudo enviar el email de verificación", "error", err)
	}
	h.setRefreshCookie(w, r, refresh)
	web.JSON(w, http.StatusCreated, res)
}

// VerifyEmail handles POST /auth/verify-email with {token} from the email link.
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	if err := h.verification.Verify(r.Context(), req.Token); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "Correo verificado"})
}

// ResendVerification handles POST /auth/resend-verification for the logged-in
// user (the user ID comes from the auth middleware).
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFrom(r.Context())
	if !ok {
		web.RespondError(w, authResource, apperror.ErrUnauthorized)
		return
	}
	if err := h.verification.Send(r.Context(), userID); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	web.JSON(w, http.StatusOK, map[string]string{"message": "Te reenviamos el enlace de verificación"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	res, refresh, err := h.svc.Login(r.Context(), req)
	if err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	h.setRefreshCookie(w, r, refresh)
	web.JSON(w, http.StatusOK, res)
}

// Refresh reads the refresh-token cookie, rotates it and returns a new access
// token. Used by the frontend when the access token expires.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshCookieName)
	if err != nil {
		web.RespondError(w, authResource, apperror.ErrUnauthorized)
		return
	}
	res, refresh, err := h.svc.Refresh(r.Context(), cookie.Value)
	if err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	h.setRefreshCookie(w, r, refresh)
	web.JSON(w, http.StatusOK, res)
}

// OAuthGoogle handles POST /auth/oauth/google with {code, codeVerifier} from
// the frontend's PKCE flow. On success it sets the refresh cookie and returns
// an access token, exactly like login.
func (h *AuthHandler) OAuthGoogle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code         string `json:"code"`
		CodeVerifier string `json:"codeVerifier"`
	}
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	res, refresh, err := h.oauth.LoginWithGoogle(r.Context(), req.Code, req.CodeVerifier)
	if err != nil {
		web.RespondError(w, authResource, err)
		return
	}
	h.setRefreshCookie(w, r, refresh)
	web.JSON(w, http.StatusOK, res)
}

// Logout revokes the refresh token and clears the cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(refreshCookieName); err == nil {
		_ = h.svc.Logout(r.Context(), cookie.Value)
	}
	h.clearRefreshCookie(w, r)
	web.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// setRefreshCookie writes the refresh token as an HttpOnly + SameSite=Strict
// cookie. Secure is set only over HTTPS so local http development still works.
func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.refreshTTL.Seconds()),
	})
}

// clearRefreshCookie expires the refresh cookie (MaxAge < 0 deletes it).
func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

package handler

import (
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

// AuthHandler exposes registration, login, refresh, logout and Google OAuth.
type AuthHandler struct {
	svc        *service.AuthService
	oauth      *service.OAuthService
	refreshTTL time.Duration
}

// NewAuthHandler wires an AuthHandler to its services and the refresh-token
// lifetime (used for the cookie's Max-Age).
func NewAuthHandler(svc *service.AuthService, oauth *service.OAuthService, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{svc: svc, oauth: oauth, refreshTTL: refreshTTL}
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
	h.setRefreshCookie(w, r, refresh)
	web.JSON(w, http.StatusCreated, res)
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

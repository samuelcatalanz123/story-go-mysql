package handler

import (
	"net/http"

	"story-go-mysql/internal/auth"
)

// Router builds the HTTP routing table. GET routes are public; write routes
// (POST/PUT/DELETE) require a valid JWT. The /auth routes are public.
func Router(
	tokens *auth.TokenManager,
	authH *AuthHandler,
	characters *CharacterHandler,
	locations *LocationHandler,
	scenes *SceneHandler,
	stories *StoryHandler,
	organizations *OrganizationHandler,
	conflicts *ConflictHandler,
	uploadDir string,
) http.Handler {
	mux := http.NewServeMux()

	// Auth (public). Refresh and logout rely on the HttpOnly refresh cookie.
	mux.HandleFunc("POST /auth/register", authH.Register)
	mux.HandleFunc("POST /auth/login", authH.Login)
	mux.HandleFunc("POST /auth/refresh", authH.Refresh)
	mux.HandleFunc("POST /auth/logout", authH.Logout)
	mux.HandleFunc("POST /auth/oauth/google", authH.OAuthGoogle)

	// Uploaded files served statically (public, read-only).
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir(uploadDir))))

	// Stories: every route is private and scoped to the authenticated user.
	mux.Handle("GET /stories", RequireAuth(tokens, http.HandlerFunc(stories.List)))
	mux.Handle("GET /stories/{id}", RequireAuth(tokens, http.HandlerFunc(stories.Get)))
	mux.Handle("POST /stories", RequireAuth(tokens, http.HandlerFunc(stories.Create)))
	mux.Handle("PUT /stories/{id}", RequireAuth(tokens, http.HandlerFunc(stories.Update)))
	mux.Handle("DELETE /stories/{id}", RequireAuth(tokens, http.HandlerFunc(stories.Delete)))

	// Characters: reads public, writes protected.
	mux.HandleFunc("GET /characters", characters.List)
	mux.HandleFunc("GET /characters/{id}", characters.Get)
	mux.Handle("POST /characters", RequireAuth(tokens, http.HandlerFunc(characters.Create)))
	mux.Handle("PUT /characters/{id}", RequireAuth(tokens, http.HandlerFunc(characters.Update)))
	mux.Handle("DELETE /characters/{id}", RequireAuth(tokens, http.HandlerFunc(characters.Delete)))
	mux.Handle("POST /characters/{id}/avatar", RequireAuth(tokens, http.HandlerFunc(characters.Avatar)))

	// Locations.
	mux.HandleFunc("GET /locations", locations.List)
	mux.HandleFunc("GET /locations/{id}", locations.Get)
	mux.Handle("POST /locations", RequireAuth(tokens, http.HandlerFunc(locations.Create)))
	mux.Handle("PUT /locations/{id}", RequireAuth(tokens, http.HandlerFunc(locations.Update)))
	mux.Handle("DELETE /locations/{id}", RequireAuth(tokens, http.HandlerFunc(locations.Delete)))
	mux.Handle("POST /locations/{id}/avatar", RequireAuth(tokens, http.HandlerFunc(locations.Avatar)))

	// Scenes.
	mux.HandleFunc("GET /scenes", scenes.List)
	mux.HandleFunc("GET /scenes/{id}", scenes.Get)
	mux.Handle("POST /scenes", RequireAuth(tokens, http.HandlerFunc(scenes.Create)))
	mux.Handle("PUT /scenes/{id}", RequireAuth(tokens, http.HandlerFunc(scenes.Update)))
	mux.Handle("DELETE /scenes/{id}", RequireAuth(tokens, http.HandlerFunc(scenes.Delete)))

	// Organizations: reads public, writes protected.
	mux.HandleFunc("GET /organizations", organizations.List)
	mux.HandleFunc("GET /organizations/{id}", organizations.Get)
	mux.Handle("POST /organizations", RequireAuth(tokens, http.HandlerFunc(organizations.Create)))
	mux.Handle("PUT /organizations/{id}", RequireAuth(tokens, http.HandlerFunc(organizations.Update)))
	mux.Handle("DELETE /organizations/{id}", RequireAuth(tokens, http.HandlerFunc(organizations.Delete)))

	// Conflicts: reads public, writes protected.
	mux.HandleFunc("GET /conflicts", conflicts.List)
	mux.HandleFunc("GET /conflicts/{id}", conflicts.Get)
	mux.Handle("POST /conflicts", RequireAuth(tokens, http.HandlerFunc(conflicts.Create)))
	mux.Handle("PUT /conflicts/{id}", RequireAuth(tokens, http.HandlerFunc(conflicts.Update)))
	mux.Handle("DELETE /conflicts/{id}", RequireAuth(tokens, http.HandlerFunc(conflicts.Delete)))

	return mux
}

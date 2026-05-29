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
) http.Handler {
	mux := http.NewServeMux()

	// Auth (public).
	mux.HandleFunc("POST /auth/register", authH.Register)
	mux.HandleFunc("POST /auth/login", authH.Login)

	// Characters: reads public, writes protected.
	mux.HandleFunc("GET /characters", characters.List)
	mux.HandleFunc("GET /characters/{id}", characters.Get)
	mux.Handle("POST /characters", RequireAuth(tokens, http.HandlerFunc(characters.Create)))
	mux.Handle("PUT /characters/{id}", RequireAuth(tokens, http.HandlerFunc(characters.Update)))
	mux.Handle("DELETE /characters/{id}", RequireAuth(tokens, http.HandlerFunc(characters.Delete)))

	// Locations.
	mux.HandleFunc("GET /locations", locations.List)
	mux.HandleFunc("GET /locations/{id}", locations.Get)
	mux.Handle("POST /locations", RequireAuth(tokens, http.HandlerFunc(locations.Create)))
	mux.Handle("PUT /locations/{id}", RequireAuth(tokens, http.HandlerFunc(locations.Update)))
	mux.Handle("DELETE /locations/{id}", RequireAuth(tokens, http.HandlerFunc(locations.Delete)))

	// Scenes.
	mux.HandleFunc("GET /scenes", scenes.List)
	mux.HandleFunc("GET /scenes/{id}", scenes.Get)
	mux.Handle("POST /scenes", RequireAuth(tokens, http.HandlerFunc(scenes.Create)))
	mux.Handle("PUT /scenes/{id}", RequireAuth(tokens, http.HandlerFunc(scenes.Update)))
	mux.Handle("DELETE /scenes/{id}", RequireAuth(tokens, http.HandlerFunc(scenes.Delete)))

	return mux
}

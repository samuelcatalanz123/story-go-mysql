package handler

import "net/http"

// Router builds the HTTP routing table from the resource handlers using
// Go's method-aware ServeMux patterns.
func Router(characters *CharacterHandler, locations *LocationHandler, scenes *SceneHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /characters", characters.Create)
	mux.HandleFunc("GET /characters", characters.List)
	mux.HandleFunc("GET /characters/{id}", characters.Get)
	mux.HandleFunc("PUT /characters/{id}", characters.Update)
	mux.HandleFunc("DELETE /characters/{id}", characters.Delete)

	mux.HandleFunc("POST /locations", locations.Create)
	mux.HandleFunc("GET /locations", locations.List)
	mux.HandleFunc("GET /locations/{id}", locations.Get)
	mux.HandleFunc("PUT /locations/{id}", locations.Update)
	mux.HandleFunc("DELETE /locations/{id}", locations.Delete)

	mux.HandleFunc("POST /scenes", scenes.Create)
	mux.HandleFunc("GET /scenes", scenes.List)
	mux.HandleFunc("GET /scenes/{id}", scenes.Get)
	mux.HandleFunc("PUT /scenes/{id}", scenes.Update)
	mux.HandleFunc("DELETE /scenes/{id}", scenes.Delete)

	return mux
}

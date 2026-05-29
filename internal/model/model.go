// Package model contains the domain types and the request DTOs shared
// across the repository, service and handler layers.
package model

import "time"

// Character is a person that can appear in scenes.
type Character struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Text      *string   `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Location is a place that can appear in scenes.
type Location struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Text      *string   `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Scene is a moment in the story that relates characters and locations.
type Scene struct {
	ID            uint64      `json:"id"`
	Title         string      `json:"title"`
	Text          *string     `json:"text"`
	StartTimeline int         `json:"startTimeline"`
	EndTimeline   int         `json:"endTimeline"`
	Characters    []Character `json:"characters"`
	Locations     []Location  `json:"locations"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

// CharacterRequest is the payload accepted when creating or updating a character.
type CharacterRequest struct {
	Title string  `json:"title"`
	Text  *string `json:"text"`
}

// LocationRequest is the payload accepted when creating or updating a location.
type LocationRequest struct {
	Title string  `json:"title"`
	Text  *string `json:"text"`
}

// SceneRequest is the payload accepted when creating or updating a scene.
type SceneRequest struct {
	Title         string   `json:"title"`
	Text          *string  `json:"text"`
	StartTimeline *int     `json:"startTimeline"`
	EndTimeline   *int     `json:"endTimeline"`
	CharacterIDs  []uint64 `json:"characterIds"`
	LocationIDs   []uint64 `json:"locationIds"`
}

// User is an authenticated account. The password hash is never serialized.
type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// RegisterRequest is the payload to create an account.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the payload to log in.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is returned by register/login: a JWT plus the user.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

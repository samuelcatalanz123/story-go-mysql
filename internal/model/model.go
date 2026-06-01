// Package model contains the domain types and the request DTOs shared
// across the repository, service and handler layers.
package model

import "time"

// Story is the top-level container that owns characters, locations and
// scenes. Each story belongs to a single user.
type Story struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Text      *string   `json:"text"`
	UserID    uint64    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// StoryRequest is the payload accepted when creating or updating a story.
type StoryRequest struct {
	Title string  `json:"title"`
	Text  *string `json:"text"`
}

// Organization is a group (guild, government, crew) that characters can
// belong to. It optionally belongs to a story.
type Organization struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Text      *string   `json:"text"`
	StoryID   *uint64   `json:"storyId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// OrganizationRequest is the payload accepted when creating or updating an
// organization.
type OrganizationRequest struct {
	Title   string  `json:"title"`
	Text    *string `json:"text"`
	StoryID *uint64 `json:"storyId"`
}

// Character is a person that can appear in scenes. A character can belong to
// several organizations (many-to-many). Organizations is omitted from the
// JSON when not loaded (e.g. characters embedded inside a scene).
type Character struct {
	ID            uint64         `json:"id"`
	Title         string         `json:"title"`
	Text          *string        `json:"text"`
	Organizations []Organization `json:"organizations,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
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

// CharacterRequest is the payload accepted when creating or updating a
// character. OrganizationIDs sets the organizations the character belongs to.
type CharacterRequest struct {
	Title           string   `json:"title"`
	Text            *string  `json:"text"`
	OrganizationIDs []uint64 `json:"organizationIds"`
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

// Page is a paginated slice of results returned by list endpoints.
type Page[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// Pagination defaults and bounds for list endpoints.
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// ListParams holds search and pagination parameters for list endpoints.
type ListParams struct {
	Query    string
	Page     int
	PageSize int
}

// Normalize clamps the params to safe values: Page >= 1 and PageSize within
// [1, MaxPageSize] (defaulting to DefaultPageSize). It returns a copy.
func (p ListParams) Normalize() ListParams {
	out := p
	if out.Page < 1 {
		out.Page = 1
	}
	if out.PageSize <= 0 {
		out.PageSize = DefaultPageSize
	}
	if out.PageSize > MaxPageSize {
		out.PageSize = MaxPageSize
	}
	return out
}

// Limit returns the SQL LIMIT (call after Normalize).
func (p ListParams) Limit() int { return p.PageSize }

// Offset returns the SQL OFFSET (call after Normalize).
func (p ListParams) Offset() int { return (p.Page - 1) * p.PageSize }

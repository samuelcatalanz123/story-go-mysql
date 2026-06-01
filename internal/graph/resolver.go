package graph

import "story-go-mysql/internal/service"

// Resolver is the dependency-injection root for GraphQL. It reuses the SAME
// services as the REST API, so both APIs share all the business logic.
type Resolver struct {
	CharacterSvc *service.CharacterService
	SceneSvc     *service.SceneService
}

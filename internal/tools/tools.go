package tools

import (
	"github.com/LackOfMorals/aura-client"
	"github.com/LackOfMorals/unofficialMcp/internal/outcomes"
)

// ToolDependencies contains all dependencies needed by tools
type ToolDependencies struct {
	AClient  *aura.AuraAPIClient
	Registry *outcomes.Registry
}

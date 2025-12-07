package implementations

import (
	"github.com/LackOfMorals/unofficialMcp/internal/outcomes"
)

// RegisterAll registers all available outcomes to the registry
func RegisterAll(registry *outcomes.Registry) error {
	// Instance management outcomes
	outcomes := []*outcomes.Outcome{
		ListInstancesOutcome(),
		GetInstanceDetailsOutcome(),
		PauseInstanceOutcome(),
		// Add more outcomes here as you implement them
		// CreateInstanceOutcome(),
		// DeleteInstanceOutcome(),
		// ResumeInstanceOutcome(),
		// ResizeInstanceOutcome(),
		// ListDatabasesOutcome(),
		// CreateDatabaseOutcome(),
		// etc.
	}

	for _, outcome := range outcomes {
		if err := registry.Register(outcome); err != nil {
			return err
		}
	}

	return nil
}

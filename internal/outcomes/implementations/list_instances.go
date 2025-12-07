package implementations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LackOfMorals/unofficialMcp/internal/outcomes"
)

// ListInstancesOutcome returns the outcome definition for listing instances
func ListInstancesOutcome() *outcomes.Outcome {
	return &outcomes.Outcome{
		ID:          "list-instances",
		Name:        "List Instances",
		Description: "Lists all Neo4j Aura DB instances in your account",
		Category:    "instance-management",
		Parameters:  []outcomes.Parameter{}, // No parameters needed
		ReadOnly:    true,
		Handler:     listInstancesHandler,
	}
}

func listInstancesHandler(ctx context.Context, deps *outcomes.OutcomeDependencies, args map[string]interface{}) (*outcomes.OutcomeResult, error) {
	// Call the Aura API to list instances
	instances, err := deps.AClient.Instances.List()
	if err != nil {
		slog.Error("Failed to list instances", "error", err)
		return &outcomes.OutcomeResult{
			Success:      false,
			Error:        fmt.Sprintf("Failed to list instances: %v", err),
			APICallsMade: []string{"GET /instances"},
		}, nil
	}

	return &outcomes.OutcomeResult{
		Success:      true,
		Data:         instances.Data,
		APICallsMade: []string{"GET /instances"},
	}, nil
}

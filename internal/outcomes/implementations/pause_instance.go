package implementations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LackOfMorals/unofficialMcp/internal/outcomes"
)

// PauseInstanceOutcome returns the outcome definition for pausing an instance
func PauseInstanceOutcome() *outcomes.Outcome {
	return &outcomes.Outcome{
		ID:          "pause-instance",
		Name:        "Pause Instance",
		Description: "Pauses a running Neo4j Aura DB instance to reduce costs while preserving data",
		Category:    "instance-management",
		Parameters: []outcomes.Parameter{
			{
				Name:        "instance_id",
				Type:        "string",
				Description: "The unique identifier of the instance to pause",
				Required:    true,
			},
		},
		ReadOnly: false,
		Handler:  pauseInstanceHandler,
	}
}

func pauseInstanceHandler(ctx context.Context, deps *outcomes.OutcomeDependencies, args map[string]interface{}) (*outcomes.OutcomeResult, error) {
	instanceID, ok := args["instance_id"].(string)
	if !ok || instanceID == "" {
		return &outcomes.OutcomeResult{
			Success: false,
			Error:   "instance_id parameter is required and must be a string",
		}, nil
	}

	apiCalls := []string{}

	// First, verify the instance exists and get its current status
	instance, err := deps.AClient.Instances.Get(instanceID)
	if err != nil {
		slog.Error("Failed to get instance", "instance_id", instanceID, "error", err)
		return &outcomes.OutcomeResult{
			Success:      false,
			Error:        fmt.Sprintf("Failed to verify instance: %v", err),
			APICallsMade: append(apiCalls, fmt.Sprintf("GET /instances/%s", instanceID)),
		}, nil
	}
	apiCalls = append(apiCalls, fmt.Sprintf("GET /instances/%s", instanceID))

	// Check if instance is already paused (example - adjust based on actual API)
	// This demonstrates how an outcome can make decisions based on current state
	if instance.Data.Status == "paused" {
		return &outcomes.OutcomeResult{
			Success: true,
			Data: map[string]interface{}{
				"message":  "Instance is already paused",
				"instance": instance.Data,
			},
			APICallsMade: apiCalls,
		}, nil
	}

	// Pause the instance
	result, err := deps.AClient.Instances.Pause(instanceID)
	if err != nil {
		slog.Error("Failed to pause instance", "instance_id", instanceID, "error", err)
		return &outcomes.OutcomeResult{
			Success:      false,
			Error:        fmt.Sprintf("Failed to pause instance: %v", err),
			APICallsMade: append(apiCalls, fmt.Sprintf("POST /instances/%s/pause", instanceID)),
		}, nil
	}
	apiCalls = append(apiCalls, fmt.Sprintf("POST /instances/%s/pause", instanceID))

	return &outcomes.OutcomeResult{
		Success:      true,
		Data:         result.Data,
		APICallsMade: apiCalls,
	}, nil
}

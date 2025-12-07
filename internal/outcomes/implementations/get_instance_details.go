package implementations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LackOfMorals/unofficialMcp/internal/outcomes"
)

// GetInstanceDetailsOutcome returns the outcome definition for getting instance details
func GetInstanceDetailsOutcome() *outcomes.Outcome {
	return &outcomes.Outcome{
		ID:          "get-instance-details",
		Name:        "Get Instance Details",
		Description: "Retrieves detailed information about a specific Neo4j Aura DB instance",
		Category:    "instance-management",
		Parameters: []outcomes.Parameter{
			{
				Name:        "instance_id",
				Type:        "string",
				Description: "The unique identifier of the instance",
				Required:    true,
			},
		},
		ReadOnly: true,
		Handler:  getInstanceDetailsHandler,
	}
}

func getInstanceDetailsHandler(ctx context.Context, deps *outcomes.OutcomeDependencies, args map[string]interface{}) (*outcomes.OutcomeResult, error) {
	instanceID, ok := args["instance_id"].(string)
	if !ok || instanceID == "" {
		return &outcomes.OutcomeResult{
			Success: false,
			Error:   "instance_id parameter is required and must be a string",
		}, nil
	}

	// Call the Aura API to get instance details
	instance, err := deps.AClient.Instances.Get(instanceID)
	if err != nil {
		slog.Error("Failed to get instance details", "instance_id", instanceID, "error", err)
		return &outcomes.OutcomeResult{
			Success:      false,
			Error:        fmt.Sprintf("Failed to get instance details: %v", err),
			APICallsMade: []string{fmt.Sprintf("GET /instances/%s", instanceID)},
		}, nil
	}

	return &outcomes.OutcomeResult{
		Success:      true,
		Data:         instance.Data,
		APICallsMade: []string{fmt.Sprintf("GET /instances/%s", instanceID)},
	}, nil
}

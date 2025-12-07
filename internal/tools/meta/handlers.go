package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/LackOfMorals/unofficialMcp/internal/outcomes"
	"github.com/LackOfMorals/unofficialMcp/internal/tools"
	"github.com/mark3labs/mcp-go/mcp"
)

// ListOutcomesHandler returns the handler for listing available outcomes
func ListOutcomesHandler(registry *outcomes.Registry) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get optional category filter
		var category string
		if catVal, ok := request.Params.Arguments["category"]; ok {
			if catStr, ok := catVal.(string); ok {
				category = catStr
			}
		}

		// Get outcomes (filtered by category if provided)
		var outcomeList []*outcomes.Outcome
		if category != "" {
			outcomeList = registry.ListByCategory(category)
		} else {
			outcomeList = registry.List()
		}

		// Build response with summary information
		type OutcomeSummary struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Category    string `json:"category"`
			ReadOnly    bool   `json:"read_only"`
		}

		summaries := make([]OutcomeSummary, 0, len(outcomeList))
		for _, outcome := range outcomeList {
			summaries = append(summaries, OutcomeSummary{
				ID:          outcome.ID,
				Name:        outcome.Name,
				Description: outcome.Description,
				Category:    outcome.Category,
				ReadOnly:    outcome.ReadOnly,
			})
		}

		// Format as JSON
		result, err := json.MarshalIndent(summaries, "", "  ")
		if err != nil {
			slog.Error("Failed to marshal outcomes", "error", err)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format outcomes: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}
}

// DescribeOutcomeHandler returns the handler for describing a specific outcome
func DescribeOutcomeHandler(registry *outcomes.Registry) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get outcome_id parameter
		outcomeID, ok := request.Params.Arguments["outcome_id"].(string)
		if !ok || outcomeID == "" {
			return mcp.NewToolResultError("outcome_id parameter is required"), nil
		}

		// Get the outcome from registry
		outcome, err := registry.Get(outcomeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Outcome not found: %s", outcomeID)), nil
		}

		// Build detailed response
		type OutcomeDetail struct {
			ID          string                `json:"id"`
			Name        string                `json:"name"`
			Description string                `json:"description"`
			Category    string                `json:"category"`
			ReadOnly    bool                  `json:"read_only"`
			Parameters  []outcomes.Parameter  `json:"parameters"`
		}

		detail := OutcomeDetail{
			ID:          outcome.ID,
			Name:        outcome.Name,
			Description: outcome.Description,
			Category:    outcome.Category,
			ReadOnly:    outcome.ReadOnly,
			Parameters:  outcome.Parameters,
		}

		// Format as JSON
		result, err := json.MarshalIndent(detail, "", "  ")
		if err != nil {
			slog.Error("Failed to marshal outcome details", "error", err)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format outcome details: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}
}

// ExecuteOutcomeHandler returns the handler for executing an outcome
func ExecuteOutcomeHandler(registry *outcomes.Registry, deps *tools.ToolDependencies) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get outcome_id parameter
		outcomeID, ok := request.Params.Arguments["outcome_id"].(string)
		if !ok || outcomeID == "" {
			return mcp.NewToolResultError("outcome_id parameter is required"), nil
		}

		// Get the outcome from registry
		outcome, err := registry.Get(outcomeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Outcome not found: %s", outcomeID)), nil
		}

		// Get arguments (if provided)
		args := make(map[string]interface{})
		if argsVal, ok := request.Params.Arguments["arguments"]; ok {
			if argsMap, ok := argsVal.(map[string]interface{}); ok {
				args = argsMap
			}
		}

		// Validate required parameters
		for _, param := range outcome.Parameters {
			if param.Required {
				if _, exists := args[param.Name]; !exists {
					return mcp.NewToolResultError(
						fmt.Sprintf("Required parameter missing: %s (%s)", param.Name, param.Description),
					), nil
				}
			}
		}

		// Execute the outcome
		result, err := outcome.Handler(ctx, deps, args)
		if err != nil {
			slog.Error("Failed to execute outcome", "outcome_id", outcomeID, "error", err)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to execute outcome: %v", err)), nil
		}

		// Format the result
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			slog.Error("Failed to marshal outcome result", "error", err)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

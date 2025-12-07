package meta

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// ListOutcomesSpec returns the tool specification for listing available outcomes
func ListOutcomesSpec() mcp.Tool {
	return mcp.NewTool("list-outcomes",
		mcp.WithDescription("Lists all available outcomes that can be achieved with the Neo4j Aura API. "+
			"An outcome represents a high-level operation (e.g., create instance, pause instance) "+
			"that may call one or more API endpoints. Use this to discover what operations are available."),
		mcp.WithTitleAnnotation("List Available Outcomes"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithString("category",
			mcp.Required(false),
			mcp.Description("Optional: Filter outcomes by category (e.g., 'instance-management', 'database-operations')")),
	)
}

// DescribeOutcomeSpec returns the tool specification for getting outcome details
func DescribeOutcomeSpec() mcp.Tool {
	return mcp.NewTool("describe-outcome",
		mcp.WithDescription("Gets detailed information about a specific outcome, including what parameters it requires, "+
			"what it does, and which API endpoints it will call. Use this before executing an outcome "+
			"to understand what inputs are needed."),
		mcp.WithTitleAnnotation("Describe Outcome"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithString("outcome_id",
			mcp.Required(true),
			mcp.Description("The unique identifier of the outcome to describe")),
	)
}

// ExecuteOutcomeSpec returns the tool specification for executing an outcome
func ExecuteOutcomeSpec() mcp.Tool {
	return mcp.NewTool("execute-outcome",
		mcp.WithDescription("Executes a specific outcome with the provided parameters. "+
			"The outcome will call the necessary API endpoints and return the results. "+
			"Use 'describe-outcome' first to understand what parameters are required."),
		mcp.WithTitleAnnotation("Execute Outcome"),
		mcp.WithReadOnlyHintAnnotation(false), // May be read-only depending on the outcome
		mcp.WithDestructiveHintAnnotation(false), // May be destructive depending on the outcome
		mcp.WithIdempotentHintAnnotation(false), // May be idempotent depending on the outcome
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithString("outcome_id",
			mcp.Required(true),
			mcp.Description("The unique identifier of the outcome to execute")),
		mcp.WithObject("arguments",
			mcp.Required(false),
			mcp.Description("Parameters for the outcome. Use 'describe-outcome' to see what parameters are required. "+
				"Format: {\"param_name\": \"param_value\", ...}")),
	)
}

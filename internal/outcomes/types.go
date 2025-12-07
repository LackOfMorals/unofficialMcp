package outcomes

import (
	"context"

	"github.com/LackOfMorals/aura-client"
)

// Outcome represents a high-level operation that can be performed
// An outcome may call one or more API endpoints to achieve its goal
type Outcome struct {
	// ID is the unique identifier for this outcome
	ID string

	// Name is a human-readable name
	Name string

	// Description explains what this outcome does
	Description string

	// Category groups related outcomes (e.g., "instance-management", "database-operations")
	Category string

	// Parameters describes what inputs this outcome needs
	Parameters []Parameter

	// ReadOnly indicates if this outcome only reads data (no mutations)
	ReadOnly bool

	// Handler is the function that executes this outcome
	Handler OutcomeHandler
}

// Parameter describes an input parameter for an outcome
type Parameter struct {
	// Name of the parameter
	Name string

	// Type of the parameter (e.g., "string", "integer", "boolean")
	Type string

	// Description of what this parameter does
	Description string

	// Required indicates if this parameter must be provided
	Required bool

	// DefaultValue is used if the parameter is not provided (only for optional parameters)
	DefaultValue interface{}

	// Enum lists valid values (if applicable)
	Enum []string
}

// OutcomeRequest contains the parameters for executing an outcome
type OutcomeRequest struct {
	// OutcomeID identifies which outcome to execute
	OutcomeID string

	// Arguments contains the parameter values
	Arguments map[string]interface{}
}

// OutcomeResult contains the result of executing an outcome
type OutcomeResult struct {
	// Success indicates if the outcome completed successfully
	Success bool

	// Data contains the result data (if successful)
	Data interface{}

	// Error contains error information (if unsuccessful)
	Error string

	// APICallsMade lists which API endpoints were called
	APICallsMade []string
}

// OutcomeDependencies contains dependencies needed by outcome handlers
type OutcomeDependencies struct {
	AClient *aura.AuraAPIClient
}

// OutcomeHandler is the function signature for outcome handlers
type OutcomeHandler func(ctx context.Context, deps *OutcomeDependencies, args map[string]interface{}) (*OutcomeResult, error)

# Outcome-Based Architecture

This MCP server uses an outcome-based architecture that provides a scalable way to expose Neo4j Aura API functionality without overwhelming the LLM with too many tools or consuming excessive tokens.

## Architecture Overview

Instead of having one tool per API endpoint, we have **three meta-tools** that work with an **outcome registry**:

1. **list-outcomes** - Lists all available outcomes
2. **describe-outcome** - Gets detailed information about a specific outcome
3. **execute-outcome** - Executes an outcome with provided parameters

### What is an Outcome?

An **outcome** represents a high-level operation or goal that can be achieved. Examples include:
- Create a Neo4j Aura instance
- Delete an instance
- Pause/resume an instance
- Get instance details
- List all instances

An outcome may call one or more API endpoints to achieve its goal. This abstraction:
- Reduces token usage (LLM only sees 3 tools instead of dozens)
- Provides better discoverability (outcomes are grouped by category)
- Allows complex operations that span multiple API calls
- Maintains clean separation of concerns

## Directory Structure

```
internal/
├── outcomes/
│   ├── types.go              # Outcome type definitions
│   ├── registry.go           # Outcome registry
│   └── implementations/
│       ├── register.go       # Registers all outcomes
│       ├── list_instances.go
│       ├── get_instance_details.go
│       └── ... (more outcome implementations)
├── tools/
│   └── meta/
│       ├── specs.go          # Meta-tool specifications
│       └── handlers.go       # Meta-tool handlers
```

## How to Add a New Outcome

### 1. Create the Outcome Implementation

Create a new file in `internal/outcomes/implementations/` (e.g., `create_instance.go`):

```go
package implementations

import (
    "context"
    "fmt"
    "github.com/LackOfMorals/unofficialMcp/internal/outcomes"
    "github.com/LackOfMorals/unofficialMcp/internal/tools"
)

func CreateInstanceOutcome() *outcomes.Outcome {
    return &outcomes.Outcome{
        ID:          "create-instance",
        Name:        "Create Instance",
        Description: "Creates a new Neo4j Aura DB instance",
        Category:    "instance-management",
        Parameters: []outcomes.Parameter{
            {
                Name:        "name",
                Type:        "string",
                Description: "Name for the new instance",
                Required:    true,
            },
            {
                Name:        "cloud_provider",
                Type:        "string",
                Description: "Cloud provider (gcp, aws, azure)",
                Required:    true,
                Enum:        []string{"gcp", "aws", "azure"},
            },
            {
                Name:        "region",
                Type:        "string",
                Description: "Cloud region",
                Required:    true,
            },
            {
                Name:        "size",
                Type:        "string",
                Description: "Instance size",
                Required:    false,
                DefaultValue: "1GB",
            },
        },
        ReadOnly: false,
        Handler:  createInstanceHandler,
    }
}

func createInstanceHandler(ctx context.Context, deps *tools.ToolDependencies, args map[string]interface{}) (*outcomes.OutcomeResult, error) {
    // Extract parameters
    name := args["name"].(string)
    cloudProvider := args["cloud_provider"].(string)
    region := args["region"].(string)
    
    size := "1GB"
    if sizeVal, ok := args["size"]; ok {
        size = sizeVal.(string)
    }

    // Call API (example - adjust based on actual Aura API client)
    instance, err := deps.AClient.Instances.Create(name, cloudProvider, region, size)
    if err != nil {
        return &outcomes.OutcomeResult{
            Success:      false,
            Error:        fmt.Sprintf("Failed to create instance: %v", err),
            APICallsMade: []string{"POST /instances"},
        }, nil
    }

    return &outcomes.OutcomeResult{
        Success:      true,
        Data:         instance,
        APICallsMade: []string{"POST /instances"},
    }, nil
}
```

### 2. Register the Outcome

Add your outcome to `internal/outcomes/implementations/register.go`:

```go
func RegisterAll(registry *outcomes.Registry) error {
    outcomes := []*outcomes.Outcome{
        ListInstancesOutcome(),
        GetInstanceDetailsOutcome(),
        CreateInstanceOutcome(),  // Add your new outcome here
        // ... other outcomes
    }

    for _, outcome := range outcomes {
        if err := registry.Register(outcome); err != nil {
            return err
        }
    }

    return nil
}
```

### 3. Test Your Outcome

Build and run the server:

```bash
go build -C cmd/neo4j-aura-mcp -o ../../bin/
./bin/neo4j-aura-mcp
```

Then test with the meta-tools:

```
# List all outcomes
list-outcomes

# Get details about your outcome
describe-outcome --outcome_id=create-instance

# Execute your outcome
execute-outcome --outcome_id=create-instance --arguments='{"name":"my-instance","cloud_provider":"gcp","region":"us-central1"}'
```

## Outcome Categories

Organize outcomes into logical categories:

- **instance-management** - Create, delete, pause, resume instances
- **database-operations** - Database-level operations
- **user-management** - User and access control
- **monitoring** - Metrics and health checks
- **configuration** - Settings and configurations

## Best Practices

1. **Keep outcomes focused** - Each outcome should do one thing well
2. **Use clear descriptions** - Help the LLM understand when to use each outcome
3. **Validate parameters** - Check required parameters and types
4. **Handle errors gracefully** - Return structured error information
5. **Track API calls** - Always list which API endpoints were called
6. **Use categories** - Group related outcomes together
7. **Document parameters** - Include type, description, and whether required

## Benefits of This Architecture

1. **Token Efficiency** - LLM only sees 3 tools initially, not dozens
2. **Better Discovery** - Outcomes can be filtered by category
3. **Scalability** - Easy to add new outcomes without changing tool structure
4. **Flexibility** - Outcomes can call multiple API endpoints if needed
5. **Type Safety** - Strong typing for parameters and results
6. **Maintainability** - Clear separation of concerns
7. **Testing** - Easy to test outcomes independently

## Example LLM Workflow

1. User: "Create a new Aura instance"
2. LLM calls: `list-outcomes` with category "instance-management"
3. LLM sees: "create-instance" outcome is available
4. LLM calls: `describe-outcome` with outcome_id "create-instance"
5. LLM learns: Needs name, cloud_provider, region parameters
6. LLM calls: `execute-outcome` with the required parameters
7. Result: Instance is created and details are returned

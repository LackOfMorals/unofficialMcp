# Import Cycle Fix

## Problem
The original code had a circular dependency:
```
outcomes/types.go → imports tools
tools/tools.go → imports outcomes
```

This caused a Go compilation error: "import cycle not allowed"

## Solution
Removed the circular dependency by:

1. **Created `OutcomeDependencies` in the `outcomes` package**
   - Moved dependency definition directly into `outcomes/types.go`
   - No longer imports `tools` package
   - Contains only what outcome handlers need: `AClient`

2. **Simplified `ToolDependencies` in the `tools` package**
   - Removed import of `outcomes` package
   - Contains only: `AClient`
   - No longer needs to reference the outcome registry

3. **Updated all outcome handlers**
   - Changed signature from `tools.ToolDependencies` to `outcomes.OutcomeDependencies`
   - Updated: `list_instances.go`, `get_instance_details.go`, `pause_instance.go`

4. **Updated meta-tool handlers**
   - Changed to accept `outcomes.OutcomeDependencies`
   - Registry passed directly as parameter, not through dependencies

## Key Changes

### Before (Circular Dependency)
```go
// outcomes/types.go
import "github.com/LackOfMorals/unofficialMcp/internal/tools"
type OutcomeHandler func(..., deps *tools.ToolDependencies, ...)

// tools/tools.go
import "github.com/LackOfMorals/unofficialMcp/internal/outcomes"
type ToolDependencies struct {
    Registry *outcomes.Registry  // ❌ Creates cycle
}
```

### After (No Cycle)
```go
// outcomes/types.go
import "github.com/LackOfMorals/aura-client"
type OutcomeDependencies struct {
    AClient *aura.AuraAPIClient
}
type OutcomeHandler func(..., deps *OutcomeDependencies, ...)

// tools/tools.go
import "github.com/LackOfMorals/aura-client"
type ToolDependencies struct {
    AClient *aura.AuraAPIClient  // ✅ No import of outcomes
}
```

## Files Modified
1. `internal/outcomes/types.go` - Added `OutcomeDependencies`, removed tools import
2. `internal/tools/tools.go` - Removed outcomes import and Registry field
3. `internal/outcomes/implementations/*.go` - Updated to use `OutcomeDependencies`
4. `internal/tools/meta/handlers.go` - Updated to use `OutcomeDependencies`
5. `internal/server/tools_register.go` - Creates `OutcomeDependencies` instead

## Build Should Now Work
```bash
go build -C cmd/neo4j-aura-mcp -o ../../bin/
```

## Note on Aura Client Methods
You may still see compilation errors about missing methods:
- `deps.AClient.Instances.Get(instanceID)` - in get_instance_details.go
- `deps.AClient.Instances.Pause(instanceID)` - in pause_instance.go

These are expected if your Aura client doesn't implement these methods yet. To fix:

**Option 1: Comment out those outcomes**
Edit `internal/outcomes/implementations/register.go`:
```go
outcomes := []*outcomes.Outcome{
    ListInstancesOutcome(),
    // GetInstanceDetailsOutcome(),  // Comment out if Get() doesn't exist
    // PauseInstanceOutcome(),        // Comment out if Pause() doesn't exist
}
```

**Option 2: Implement the missing methods in your Aura client**

**Option 3: Adjust the code to use existing Aura client methods**

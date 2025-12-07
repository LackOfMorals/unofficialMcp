# Architecture Comparison: Old vs New

## Visual Overview

### OLD: One Tool Per Endpoint
```
┌─────────────────────────────────────────┐
│           MCP Server                     │
├─────────────────────────────────────────┤
│  Tool: list-instances                    │
│  Tool: get-instance                      │
│  Tool: create-instance                   │
│  Tool: delete-instance                   │
│  Tool: pause-instance                    │
│  Tool: resume-instance                   │
│  Tool: resize-instance                   │
│  Tool: list-databases                    │
│  Tool: create-database                   │
│  Tool: delete-database                   │
│  ... (potentially 50+ tools)             │
└─────────────────────────────────────────┘
         ↓ All tools exposed to LLM
┌─────────────────────────────────────────┐
│  LLM sees ALL tools in every request    │
│  High token usage                        │
│  Overwhelming choice                     │
└─────────────────────────────────────────┘
```

### NEW: Three Meta-Tools + Outcome Registry
```
┌─────────────────────────────────────────┐
│           MCP Server                     │
├─────────────────────────────────────────┤
│  Tool: list-outcomes ────────┐          │
│  Tool: describe-outcome ─────┼──────┐   │
│  Tool: execute-outcome ──────┘      │   │
└─────────────────────────────────────┼───┘
                                      │
         ┌────────────────────────────┘
         ↓ Only 3 tools exposed to LLM
┌─────────────────────────────────────────┐
│  LLM sees only 3 tools initially         │
│  Low token usage                         │
│  Progressive discovery                   │
└─────────────────────────────────────────┘
         ↓ When needed, queries registry
┌─────────────────────────────────────────┐
│        Outcome Registry                  │
├─────────────────────────────────────────┤
│  • list-instances                        │
│  • get-instance-details                  │
│  • create-instance                       │
│  • delete-instance                       │
│  • pause-instance                        │
│  • resume-instance                       │
│  ... (unlimited outcomes)                │
└─────────────────────────────────────────┘
```

## Code Comparison

### OLD: Direct Tool Registration

**tools_register.go:**
```go
func (s *Neo4jMCPServer) RegisterTools() error {
    all := []server.ServerTool{
        {Tool: ListInstanceSpec(), Handler: ListInstanceHandler(deps)},
        {Tool: GetInstanceSpec(), Handler: GetInstanceHandler(deps)},
        {Tool: CreateInstanceSpec(), Handler: CreateInstanceHandler(deps)},
        {Tool: DeleteInstanceSpec(), Handler: DeleteInstanceHandler(deps)},
        // ... 50+ more tools
    }
    s.MCPServer.AddTools(all...)
}
```

**Each tool needs:**
- Separate spec file
- Separate handler file
- Registration in tools_register.go

### NEW: Outcome-Based Registration

**tools_register.go:**
```go
func (s *Neo4jMCPServer) RegisterTools() error {
    registry := outcomes.NewRegistry()
    implementations.RegisterAll(registry)  // Registers all outcomes
    
    deps := &tools.ToolDependencies{
        AClient:  s.aClient,
        Registry: registry,
    }
    
    metaTools := []server.ServerTool{
        {Tool: meta.ListOutcomesSpec(), Handler: meta.ListOutcomesHandler(registry)},
        {Tool: meta.DescribeOutcomeSpec(), Handler: meta.DescribeOutcomeHandler(registry)},
        {Tool: meta.ExecuteOutcomeSpec(), Handler: meta.ExecuteOutcomeHandler(registry, deps)},
    }
    
    s.MCPServer.AddTools(metaTools...)
}
```

**Each outcome needs:**
- Single implementation file
- Registration in register.go (one line)

## Token Usage Comparison

### Scenario: LLM wants to list instances

**OLD (One Tool Per Endpoint):**
```
Context sent to LLM on EVERY request:
- list-instances tool spec (50 tokens)
- get-instance tool spec (50 tokens)
- create-instance tool spec (100 tokens)
- delete-instance tool spec (80 tokens)
- pause-instance tool spec (80 tokens)
... (50+ more tool specs)
= ~3,000+ tokens per request
```

**NEW (Outcome-Based):**
```
Initial request context:
- list-outcomes tool spec (40 tokens)
- describe-outcome tool spec (40 tokens)
- execute-outcome tool spec (50 tokens)
= ~130 tokens per request

Only when needed:
- LLM calls list-outcomes (100 tokens response)
- LLM calls execute-outcome (50 tokens response)
= ~280 tokens total
```

**Savings: ~90% reduction in tool-related tokens**

## LLM Interaction Flow

### OLD Approach
```
User: "List my instances"
  ↓
LLM sees 50+ tools
  ↓
LLM picks "list-instances" tool
  ↓
Result returned
```
Simple but doesn't scale.

### NEW Approach
```
User: "List my instances"
  ↓
LLM sees 3 tools
  ↓
LLM calls "list-outcomes" (discovers operations)
  ↓
LLM identifies "list-instances" outcome
  ↓
LLM calls "execute-outcome" with outcome_id
  ↓
Result returned
```
One extra step but much more scalable.

## Scalability Comparison

| Metric | OLD | NEW |
|--------|-----|-----|
| **Tools exposed to LLM** | N (50+) | 3 |
| **Token usage per request** | High (~3000) | Low (~130) |
| **Complexity to add operation** | 3 files + registration | 1 file + registration |
| **Discoverability** | Poor (too many) | Good (categorized) |
| **Multi-endpoint operations** | Need multiple tools | Single outcome |
| **Context window used** | High | Low |

## Migration Path

### Phase 1: Add New System (Current)
- ✅ Create outcome infrastructure
- ✅ Create three meta-tools
- ✅ Implement example outcomes
- ⏳ Test and validate

### Phase 2: Expand Outcomes (Next)
- Add more outcome implementations
- Cover all Aura API endpoints
- Organize by category

### Phase 3: Remove Old System (Future)
- Delete `internal/tools/auraapi/` directory
- Remove old tool registrations
- Update documentation

## Key Benefits Recap

1. **Scalability**: Can add hundreds of operations without overwhelming LLM
2. **Efficiency**: 90% reduction in tool-related token usage
3. **Flexibility**: Outcomes can combine multiple API calls
4. **Organization**: Categories help with discovery
5. **Maintenance**: One file per operation instead of three
6. **Evolution**: Easy to add new operations without changing tool structure

## When to Use Each Pattern

**Use OLD (Direct Tools) when:**
- You have < 10 operations total
- Operations are simple (1 API call each)
- Token usage is not a concern

**Use NEW (Outcome-Based) when:**
- You have 10+ operations (or plan to)
- Operations may be complex (multiple API calls)
- Token efficiency matters
- You want better organization and scalability

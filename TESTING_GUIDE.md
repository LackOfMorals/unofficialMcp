# Testing Guide for Outcome-Based Architecture

## Quick Start

### 1. Build the Server
```bash
cd /Users/jgiffard/Projects/unofficialMcp
go build -C cmd/neo4j-aura-mcp -o ../../bin/
```

If you encounter any compilation errors, they're likely due to:
- Methods on the Aura API client that don't exist yet (Get, Pause, etc.)
- We've used placeholder methods that you'll need to adjust based on your actual Aura client library

### 2. Run the Server
```bash
./bin/neo4j-aura-mcp
```

The server will start and listen for MCP protocol messages via stdio.

## Testing with Claude

Once the server is running, you can test it through Claude. The server will now expose three tools:

### Test 1: List Available Outcomes
Ask Claude:
```
"Can you show me what operations are available?"
```

Claude should call `list-outcomes` and see:
```json
[
  {
    "id": "list-instances",
    "name": "List Instances",
    "description": "Lists all Neo4j Aura DB instances in your account",
    "category": "instance-management",
    "read_only": true
  },
  {
    "id": "get-instance-details",
    "name": "Get Instance Details",
    "description": "Retrieves detailed information about a specific Neo4j Aura DB instance",
    "category": "instance-management",
    "read_only": true
  },
  {
    "id": "pause-instance",
    "name": "Pause Instance",
    "description": "Pauses a running Neo4j Aura DB instance to reduce costs while preserving data",
    "category": "instance-management",
    "read_only": false
  }
]
```

### Test 2: Describe an Outcome
Ask Claude:
```
"What parameters do I need to get instance details?"
```

Claude should call `describe-outcome` with `outcome_id: "get-instance-details"` and see:
```json
{
  "id": "get-instance-details",
  "name": "Get Instance Details",
  "description": "Retrieves detailed information about a specific Neo4j Aura DB instance",
  "category": "instance-management",
  "read_only": true,
  "parameters": [
    {
      "Name": "instance_id",
      "Type": "string",
      "Description": "The unique identifier of the instance",
      "Required": true,
      "DefaultValue": null,
      "Enum": null
    }
  ]
}
```

### Test 3: Execute an Outcome
Ask Claude:
```
"Can you list all my Aura instances?"
```

Claude should call `execute-outcome` with `outcome_id: "list-instances"` and get your actual instances back.

### Test 4: Execute with Parameters
Ask Claude:
```
"Can you get the details for instance 07a00209?"
```

Claude should:
1. Possibly call `describe-outcome` to understand what parameters are needed
2. Call `execute-outcome` with:
   - `outcome_id: "get-instance-details"`
   - `arguments: {"instance_id": "07a00209"}`

## Debugging

### Check Server Logs
The server logs to stderr. Look for:
- Tool registration messages
- Incoming tool calls
- Outcome execution results
- Any errors

### Common Issues

**Issue**: Compilation errors about missing methods
- **Fix**: The Aura client might not have all the methods we're calling (Get, Pause, etc.)
- **Solution**: Comment out the outcomes that use missing methods in `register.go`, or implement them in the Aura client

**Issue**: Outcomes not appearing in list-outcomes
- **Fix**: Check that the outcome is registered in `implementations/register.go`
- **Solution**: Add it to the `outcomes` slice

**Issue**: Parameter validation errors
- **Fix**: The arguments passed to execute-outcome don't match the parameter definitions
- **Solution**: Use describe-outcome to see what parameters are required

## Verifying the Implementation

### Minimal Test (No Aura API Calls)

To test the meta-tool infrastructure without making actual API calls:

1. List outcomes should always work (just reads from registry)
2. Describe outcome should work for any registered outcome
3. Execute outcome will fail if the Aura API client methods don't exist

### Full Test (With Aura API)

Once the Aura client methods are confirmed:

1. List instances should return your actual instances
2. Get instance details should return details for a specific instance
3. Pause instance should pause an instance (only if the API supports it)

## Expected LLM Behavior

The LLM should now follow this pattern:

1. **Discovery**: Use `list-outcomes` to see what's available
2. **Learning**: Use `describe-outcome` to understand parameters
3. **Execution**: Use `execute-outcome` to perform the action

This is more token-efficient because:
- The LLM only loads outcome details when needed
- Not all outcomes are in context at once
- The LLM can filter by category to find relevant operations

## Next Steps After Testing

Once basic functionality is confirmed:

1. **Add more outcomes** - See the comments in `register.go` for suggestions
2. **Adjust Aura API calls** - Make sure the client methods match your actual Aura client library
3. **Test read-only mode** - Verify that read-only configuration works as expected
4. **Add error handling** - Improve error messages for better LLM understanding

## Troubleshooting Checklist

- [ ] Server compiles without errors
- [ ] Server starts and listens for input
- [ ] list-outcomes returns the three registered outcomes
- [ ] describe-outcome returns parameter details for each outcome
- [ ] execute-outcome works for list-instances
- [ ] Aura API client methods exist and work correctly
- [ ] Error messages are clear and actionable

If you get stuck, check:
1. The Aura client library implementation
2. The go.mod dependencies
3. The server logs for specific error messages

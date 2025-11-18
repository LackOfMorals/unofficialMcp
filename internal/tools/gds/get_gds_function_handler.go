package gds

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/neo4j/mcp/internal/analytics"
	"github.com/neo4j/mcp/internal/database"
	"github.com/neo4j/mcp/internal/tools"
)

const GetGdsProceduresQuery = `
CALL gds.list() YIELD name, description, signature, type
WHERE type = "procedure"
AND name CONTAINS "stream"
AND NOT (name CONTAINS "estimate")
RETURN name, description, signature, type`

func GetGDSFunctionDetailHandler(deps *tools.ToolDependencies) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetGdsProcedure(ctx, request, deps.DBService, deps.AnalyticsService)
	}
}

func handleGetGdsProcedure(ctx context.Context, request mcp.CallToolRequest, dbService database.Service, asService analytics.Service) (*mcp.CallToolResult, error) {
	if dbService == nil {
		errMessage := "Database service is not initialized"
		log.Printf("%s", errMessage)
		return mcp.NewToolResultError(errMessage), nil
	}

	if asService == nil {
		errMessage := "Analytics service is not initialized"
		log.Printf("%s", errMessage)
		return mcp.NewToolResultError(errMessage), nil
	}
	asService.EmitEvent(asService.NewToolsEvent("get-gds-function-detail"))

	var args GetGDSFunctionDetailsInput
	// Bind arguments to the struct
	if err := request.BindArguments(&args); err != nil {
		log.Printf("Error binding arguments: %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Gets the list of available GDS functions
	Query := fmt.Sprintf(`SHOW PROCEDURES YIELD * WHERE name = '%s' RETURN name, description, signature, argumentDescription, returnDescription`, args.FunctionName)

	// Execute the Cypher query using the database service (now confirmed read-only)
	records, err := dbService.ExecuteReadQuery(ctx, Query, nil)
	if err != nil {
		log.Printf("Error executing Cypher query: %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	response, err := dbService.Neo4jRecordsToJSON(records)
	if err != nil {
		log.Printf("Error formatting query results: %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(response), nil
}

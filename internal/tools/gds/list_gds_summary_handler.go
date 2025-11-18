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

const listGdsSummaryProceduresQuery = `
CALL gds.list() YIELD name, description,type
WHERE type = "procedure"
AND name CONTAINS "stream"
AND NOT (name CONTAINS "estimate")
RETURN name, description`

func ListSummaryGDSProceduresHandler(deps *tools.ToolDependencies) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListSummaryGdsProcedures(ctx, deps.DBService, deps.AnalyticsService)
	}
}

func handleListSummaryGdsProcedures(ctx context.Context, dbService database.Service, asService analytics.Service) (*mcp.CallToolResult, error) {
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
	asService.EmitEvent(asService.NewToolsEvent("list-gds-summary-procedures"))

	records, err := dbService.ExecuteReadQuery(ctx, listGdsSummaryProceduresQuery, nil)
	if err != nil {
		formattedErrorMessage := fmt.Errorf("failed to execute query listGdsSummaryProceduresQuery : %v. Ensure that the Graph Data Science (GDS) library is installed and properly configured in your Neo4j database", err)
		log.Printf("%s", formattedErrorMessage.Error())
		return mcp.NewToolResultError(formattedErrorMessage.Error()), nil
	}

	response, err := dbService.Neo4jRecordsToJSON(records)
	if err != nil {
		log.Printf("Failed to format query results to JSON: %v", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(response), nil
}

package gds

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func ListSummaryGDSProceduresSpec() mcp.Tool {
	return mcp.NewTool("list-gds-summary-procedures",
		mcp.WithDescription(
			"Graph science and analytics functions help you with centrality, community detection, similarity, path finding, and identifying dependencies between nodes. "+
				"Use this tool to discover what graph science and analytics functions are available in the current Neo4j environment. "+
				"It returns a structured list with the name of each function and what it does."+
				"Give the name of a GDS function to get-gds-function-detail tool to get the full description."+
				"Do this before any reasoning, query generation, or analysis so you know what capabilities exist."),
		mcp.WithTitleAnnotation("A summary of available Neo4j GDS procedures"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(true),
	)
}

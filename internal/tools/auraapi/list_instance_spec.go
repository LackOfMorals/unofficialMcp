package auraapi

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func ListInstanceSpec() mcp.Tool {
	return mcp.NewTool("list-instances",
		mcp.WithDescription("lists all instances in Neo4j Aura DB cloud platform"),
		mcp.WithTitleAnnotation("List instances"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(true),
	)
}

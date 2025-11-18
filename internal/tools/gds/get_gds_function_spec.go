package gds

import "github.com/mark3labs/mcp-go/mcp"

type GetGDSFunctionDetailsInput struct {
	FunctionName string `json:"functionName" jsonschema:"description=Gets the details of a named GDS function."`
}

func GetGDSFunctionDetailSpec() mcp.Tool {
	return mcp.NewTool("get-gds-function-detail",
		mcp.WithDescription("Before you run a Neo4j Graph Data Science function, use this tool to obtain  the full details for it including parameters to use and its reponse."),
		mcp.WithTitleAnnotation("Gets the detalis for a GDS function"),
		mcp.WithInputSchema[GetGDSFunctionDetailsInput](),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(true),
	)
}

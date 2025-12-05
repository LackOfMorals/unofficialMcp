package aura

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func CreateInstanceHandler() func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleCreateInstance(ctx, request)
	}
}

func handleCreateInstance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	return mcp.NewToolResultText(""), nil
}

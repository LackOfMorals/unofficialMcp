package auraapi

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/LackOfMorals/unofficialMcp/internal/tools"
	"github.com/mark3labs/mcp-go/mcp"
)

func ListInstanceHandler(deps *tools.ToolDependencies) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListInstance(ctx, deps)
	}
}

func handleListInstance(ctx context.Context, deps *tools.ToolDependencies) (*mcp.CallToolResult, error) {

	// Get the list of instances from aura
	instances, err := deps.AClient.Instances.List()
	if err != nil {
		slog.Error("Failed to list instances: %v", slog.Any("Error: %v", err))
	}

	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "%#v", instances.Data)

	return mcp.NewToolResultText(buffer.String()), nil
}

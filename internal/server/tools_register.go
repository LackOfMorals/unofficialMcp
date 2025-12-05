package server

import (
	"github.com/LackOfMorals/unofficialMcp/internal/tools"
	"github.com/LackOfMorals/unofficialMcp/internal/tools/auraapi"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all enabled MCP tools and adds them to the provided MCP server.
// Tools are filtered according to the server configuration. For example, when the read-only
// mode is enabled (e.g. via the NEO4J_READ_ONLY environment variable or the Config.ReadOnly flag),
// any tool that performs state mutation will be excluded; only tools annotated as read-only will be registered.
// Note: this read-only filtering relies on the tool annotation "readonly" (ReadOnlyHint). If the annotation
// is not defined or is set to false, the tool will be added (i.e., only tools with readonly=true are filtered in read-only mode).
func (s *Neo4jMCPServer) RegisterTools() error {
	deps := &tools.ToolDependencies{
		AClient: s.aClient,
	}

	all := getAllTools(deps)

	// If read-only mode is enabled, expose only tools annotated as read-only.
	if s.config != nil && s.config.ReadOnly == "true" {
		readOnlyTools := make([]server.ServerTool, 0, len(all))
		for _, t := range all {
			if t.Tool.Annotations.ReadOnlyHint != nil && *t.Tool.Annotations.ReadOnlyHint {
				readOnlyTools = append(readOnlyTools, t)
			}
		}
		s.MCPServer.AddTools(readOnlyTools...)
		return nil
	}

	s.MCPServer.AddTools(all...)
	return nil
}

// getAllTools returns all available tools with their specs and handlers
func getAllTools(deps *tools.ToolDependencies) []server.ServerTool {
	return []server.ServerTool{
		// Aura infrastructure management
		{
			Tool:    auraapi.ListInstanceSpec(),
			Handler: auraapi.ListInstanceHandler(deps),
		},
	}
}

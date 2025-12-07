package server

import (
	"github.com/LackOfMorals/unofficialMcp/internal/outcomes"
	"github.com/LackOfMorals/unofficialMcp/internal/outcomes/implementations"
	"github.com/LackOfMorals/unofficialMcp/internal/tools/meta"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers the three meta-tools that provide access to all outcomes
// These meta-tools are:
// 1. list-outcomes: Lists all available outcomes
// 2. describe-outcome: Gets details about a specific outcome
// 3. execute-outcome: Executes an outcome with provided parameters
func (s *Neo4jMCPServer) RegisterTools() error {
	// Create and populate the outcome registry
	registry := outcomes.NewRegistry()
	if err := implementations.RegisterAll(registry); err != nil {
		return err
	}

	// Create dependencies for outcome handlers
	deps := &outcomes.OutcomeDependencies{
		AClient: s.aClient,
	}

	// Register the three meta-tools
	metaTools := []server.ServerTool{
		{
			Tool:    meta.ListOutcomesSpec(),
			Handler: meta.ListOutcomesHandler(registry),
		},
		{
			Tool:    meta.DescribeOutcomeSpec(),
			Handler: meta.DescribeOutcomeHandler(registry),
		},
		{
			Tool:    meta.ExecuteOutcomeSpec(),
			Handler: meta.ExecuteOutcomeHandler(registry, deps),
		},
	}

	// In read-only mode, we could filter outcomes but the execute-outcome tool
	// will enforce read-only restrictions based on the outcome's ReadOnly flag
	// For now, we register all three tools regardless of read-only mode
	// The execute-outcome handler can check s.config.ReadOnly if needed

	s.MCPServer.AddTools(metaTools...)
	return nil
}

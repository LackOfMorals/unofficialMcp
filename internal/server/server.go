package server

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LackOfMorals/aura-client"
	"github.com/LackOfMorals/unofficialMcp/internal/config"
	"github.com/mark3labs/mcp-go/server"
)

// Neo4jMCPServer represents the MCP server instance
type Neo4jMCPServer struct {
	MCPServer *server.MCPServer
	config    *config.Config
	aClient   *aura.AuraAPIClient
	version   string
	logger    *slog.Logger
}

// NewNeo4jMCPServer creates a new MCP server instance
// The config parameter is expected to be already validated
func NewNeo4jMCPServer(version string, cfg *config.Config) *Neo4jMCPServer {

	// Initialize Aura API Client
	client, err := aura.NewClient(
		aura.WithCredentials(cfg.ClientId, cfg.ClientSecret),
	)
	if err != nil {
		slog.Error("Failed to create Aura API client: %v", slog.Any("Error: %v", err))
		os.Exit(1)
	}

	mcpServer := server.NewMCPServer(
		"neo4j-aura-mcp",
		version,
		server.WithToolCapabilities(true),
		server.WithInstructions("This MCP Server provides tool calling to manage your Neo4j Aura Infrastructure"),
	)
	return &Neo4jMCPServer{
		MCPServer: mcpServer,
		config:    cfg,
		aClient:   client,
		version:   version,
	}
}

// Start initializes and starts the MCP server using stdio transport
func (s *Neo4jMCPServer) Start() error {
	slog.Info("Starting Aura infrastructure MCP Server...")

	// Register tools
	if err := s.RegisterTools(); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}
	slog.Info("Started Aura infrastructure MCP Server. Now listening for input...")
	// Note: ServeStdio handles its own signal management for graceful shutdown
	return server.ServeStdio(s.MCPServer)
}

// Stop gracefully stops the server
func (s *Neo4jMCPServer) Stop() error {
	slog.Info("Stopping Aura infrastructure MCP Server...")
	// Currently no cleanup needed - the MCP server handles its own lifecycle
	// Database service cleanup is handled by the caller (main.go)
	return nil
}

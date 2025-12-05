package main

import (
	"log/slog"
	"os"

	"github.com/LackOfMorals/unofficialMcp/internal/cli"
	"github.com/LackOfMorals/unofficialMcp/internal/config"
	"github.com/LackOfMorals/unofficialMcp/internal/server"
)

// go build -C cmd/neo4j-aura-mcp -o ../../bin/ -ldflags "-X 'main.Version=001' "
var Version = "development"

func main() {
	// Enable debug-level logging to stderr
	opts := &slog.HandlerOptions{Level: slog.LevelWarn}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Handle CLI arguments (version, help, etc.)
	cli.HandleArgs(Version)

	// get config from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration: %v", slog.Any("Error: %v", err))
		os.Exit(1)
	}

	// Create and configure the MCP server
	mcpServer := server.NewNeo4jMCPServer(Version, cfg)

	// Gracefully handle shutdown
	defer func() {
		if err := mcpServer.Stop(); err != nil {
		}
	}()

	// Start the server (this blocks until the server is stopped)
	if err := mcpServer.Start(); err != nil {
		slog.Error("Server error: %v", slog.Any("Error: %v", err))
		return // so that defer can run
	}

}

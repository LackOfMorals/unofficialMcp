package main

import (
	"log"

	"github.com/LackOfMorals/unofficialMcp/internal/cli"
	"github.com/LackOfMorals/unofficialMcp/internal/config"
	"github.com/LackOfMorals/unofficialMcp/internal/server"
)

// go build -C cmd/neo4j-aura-mcp -o ../../bin/ -ldflags "-X 'main.Version=001' "
var Version = "development"

func main() {
	// Handle CLI arguments (version, help, etc.)
	cli.HandleArgs(Version)

	// get config from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and configure the MCP server
	mcpServer := server.NewNeo4jMCPServer(Version, cfg)

	// Gracefully handle shutdown
	defer func() {
		if err := mcpServer.Stop(); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
	}()

	// Start the server (this blocks until the server is stopped)
	if err := mcpServer.Start(); err != nil {
		log.Printf("Server error: %v", err)
		return // so that defer can run
	}

}

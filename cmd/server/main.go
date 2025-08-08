package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"

	"ez-web-search/internal/config"
	"ez-web-search/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create a new MCP server
	s := server.NewMCPServer(
		cfg.Server.Name,
		cfg.Server.Version,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Create MCP handler
	mcpHandler := handlers.NewMCPHandler(cfg)

	// Add web search tool
	webSearchTool := mcpHandler.GetWebSearchTool()
	s.AddTool(webSearchTool, mcpHandler.HandleWebSearch)

	// Add web fetch tool
	webFetchTool := mcpHandler.GetWebFetchTool()
	s.AddTool(webFetchTool, mcpHandler.HandleWebFetch)

	// Add ping tool
	pingTool := mcpHandler.GetPingTool()
	s.AddTool(pingTool, mcpHandler.HandlePing)

	// Start the stdio server
	log.Printf("Starting %s v%s...", cfg.Server.Name, cfg.Server.Version)
	log.Println("Configuration:")
	log.Printf("  - BigModel Token: %s...%s", cfg.BigModel.Token[:8], cfg.BigModel.Token[len(cfg.BigModel.Token)-8:])
	log.Printf("  - Web Fetch Timeout: %v", cfg.WebFetch.Timeout)
	log.Printf("  - User Agent Rotation: %v", cfg.WebFetch.UserAgentRotate)
	log.Printf("  - Max Content Size: %d", cfg.WebFetch.MaxContentSize)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

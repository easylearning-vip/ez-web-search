package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"ez-web-search/internal/config"
	"ez-web-search/internal/handlers"
)

func main() {
	// Parse command line flags
	var token string
	var showHelp bool
	flag.StringVar(&token, "token", "", "BigModel API token (overrides environment variable)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.Parse()

	if showHelp {
		showUsage()
		return
	}

	// Load configuration
	cfg := config.Load()

	// Override token if provided via command line
	if token != "" {
		cfg.BigModel.Token = token
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

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
	if cfg.BigModel.Token != "" {
		log.Printf("  - BigModel Token: %s...%s", cfg.BigModel.Token[:8], cfg.BigModel.Token[len(cfg.BigModel.Token)-8:])
	} else {
		log.Printf("  - BigModel Token: Not configured (using default)")
	}
	log.Printf("  - Web Fetch Timeout: %v", cfg.WebFetch.Timeout)
	log.Printf("  - User Agent Rotation: %v", cfg.WebFetch.UserAgentRotate)
	log.Printf("  - Max Content Size: %d", cfg.WebFetch.MaxContentSize)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func showUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("  -token string     BigModel API token (overrides environment variable)")
	fmt.Println("  -help             Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  BIGMODEL_TOKEN    BigModel API token (required)")
	fmt.Println("  BIGMODEL_BASE_URL BigModel API base URL (default: https://open.bigmodel.cn/api/paas/v4/web_search)")
	fmt.Println("  BIGMODEL_TIMEOUT  BigModel API timeout (default: 30s)")
	fmt.Println("  WEBFETCH_TIMEOUT  Web fetch timeout (default: 30s)")
	fmt.Println("  WEBFETCH_MAX_CONTENT_SIZE Maximum content size to fetch (default: 5000)")
	fmt.Println("  WEBFETCH_USER_AGENT_ROTATE Enable user agent rotation (default: true)")
}

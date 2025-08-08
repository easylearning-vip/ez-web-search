package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"ez-web-search/internal/config"
	"ez-web-search/internal/services"
	"ez-web-search/pkg/types"
)

// MCPHandler handles MCP tool requests
type MCPHandler struct {
	config           *config.Config
	webSearchService *services.WebSearchService
	webFetchService  *services.WebFetchService
}

// NewMCPHandler creates a new MCP handler
func NewMCPHandler(cfg *config.Config) *MCPHandler {
	return &MCPHandler{
		config:           cfg,
		webSearchService: services.NewWebSearchService(cfg),
		webFetchService:  services.NewWebFetchService(cfg),
	}
}

// HandleWebSearch handles web search tool requests
func (h *MCPHandler) HandleWebSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract query parameter
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid query parameter: %v", err)), nil
	}

	// Extract search_engine parameter (optional, defaults to config default)
	searchEngine := h.config.BigModel.SearchEngine
	if engineVal, exists := request.GetArguments()["search_engine"]; exists {
		if strVal, ok := engineVal.(string); ok && strVal != "" {
			// Validate search engine
			validEngines := map[string]bool{
				"search_std":       true,
				"search_pro":       true,
				"search_pro_sogou": true,
				"search_pro_quark": true,
			}
			if validEngines[strVal] {
				searchEngine = strVal
			}
		}
	}

	// Extract search_intent parameter (optional, defaults to false)
	searchIntent := false
	if intentVal, exists := request.GetArguments()["search_intent"]; exists {
		if boolVal, ok := intentVal.(bool); ok {
			searchIntent = boolVal
		}
	}

	// Perform the search
	opts := types.WebSearchOptions{
		Query:        query,
		SearchEngine: searchEngine,
		SearchIntent: searchIntent,
	}

	searchResp, err := h.webSearchService.Search(ctx, opts)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	// Format the response
	resultText := h.webSearchService.FormatSearchResponse(searchResp, query, searchEngine)
	return mcp.NewToolResultText(resultText), nil
}

// HandleWebFetch handles web fetch tool requests
func (h *MCPHandler) HandleWebFetch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract URL parameter
	targetURL, err := request.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid url parameter: %v", err)), nil
	}

	// Extract optional parameters
	includeLinks := false
	includeImages := false

	if linksVal, exists := request.GetArguments()["include_links"]; exists {
		if boolVal, ok := linksVal.(bool); ok {
			includeLinks = boolVal
		}
	}

	if imagesVal, exists := request.GetArguments()["include_images"]; exists {
		if boolVal, ok := imagesVal.(bool); ok {
			includeImages = boolVal
		}
	}

	// Fetch the web page
	opts := types.WebFetchOptions{
		URL:           targetURL,
		IncludeLinks:  includeLinks,
		IncludeImages: includeImages,
	}

	content, err := h.webFetchService.FetchWebPage(ctx, opts)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch web page: %v", err)), nil
	}

	// Format the response
	resultText := h.webFetchService.FormatWebPageContent(content, includeLinks, includeImages)
	return mcp.NewToolResultText(resultText), nil
}

// HandlePing handles ping tool requests
func (h *MCPHandler) HandlePing(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("pong"), nil
}

// GetWebSearchTool returns the web search tool definition
func (h *MCPHandler) GetWebSearchTool() mcp.Tool {
	return mcp.NewTool("web_search",
		mcp.WithDescription("Search the web using BigModel Web Search API with configurable search engines"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query to execute"),
		),
		mcp.WithString("search_engine",
			mcp.Description("Search engine to use: search_std (default), search_pro, search_pro_sogou, search_pro_quark"),
		),
		mcp.WithBoolean("search_intent",
			mcp.Description("Whether to enable search intent analysis (default: false)"),
		),
	)
}

// GetWebFetchTool returns the web fetch tool definition
func (h *MCPHandler) GetWebFetchTool() mcp.Tool {
	return mcp.NewTool("web_fetch",
		mcp.WithDescription("Fetch and extract content from a web page with anti-bot protection"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL of the web page to fetch"),
		),
		mcp.WithBoolean("include_links",
			mcp.Description("Whether to include extracted links (default: false)"),
		),
		mcp.WithBoolean("include_images",
			mcp.Description("Whether to include extracted images (default: false)"),
		),
	)
}

// GetPingTool returns the ping tool definition
func (h *MCPHandler) GetPingTool() mcp.Tool {
	return mcp.NewTool("ping",
		mcp.WithDescription("Simple ping tool to test server connectivity"),
	)
}

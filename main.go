package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	// BigModel Web Search API endpoint
	webSearchURL = "https://open.bigmodel.cn/api/paas/v4/web_search"
	// Default API token (can be overridden via environment variable)
	defaultToken = "0f405f7a11b946298b154f042a70f12b.s6VO3ITALpa3bhDo"
)

// WebSearchRequest represents the request structure for BigModel Web Search API
type WebSearchRequest struct {
	SearchQuery  string `json:"search_query"`
	SearchEngine string `json:"search_engine"`
	SearchIntent bool   `json:"search_intent"`
}

// WebSearchResponse represents the response structure from BigModel Web Search API
type WebSearchResponse struct {
	ID           string         `json:"id"`
	Created      int64          `json:"created"`
	RequestID    string         `json:"request_id"`
	SearchIntent []SearchIntent `json:"search_intent"`
	SearchResult []SearchResult `json:"search_result"`
}

// SearchIntent represents search intent information
type SearchIntent struct {
	Query    string `json:"query"`
	Intent   string `json:"intent"`
	Keywords string `json:"keywords"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Link        string `json:"link"`
	Media       string `json:"media"`
	Icon        string `json:"icon"`
	Refer       string `json:"refer"`
	PublishDate string `json:"publish_date"`
}

// WebPageContent represents the content of a fetched web page
type WebPageContent struct {
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Description string            `json:"description"`
	Keywords    string            `json:"keywords"`
	Author      string            `json:"author"`
	Language    string            `json:"language"`
	Headers     map[string]string `json:"headers"`
	Links       []string          `json:"links"`
	Images      []string          `json:"images"`
	StatusCode  int               `json:"status_code"`
	ContentType string            `json:"content_type"`
}

// WebSearchClient handles communication with BigModel Web Search API
type WebSearchClient struct {
	token      string
	httpClient *http.Client
}

// WebFetchClient handles web page fetching and content extraction
type WebFetchClient struct {
	httpClient *http.Client
}

// NewWebSearchClient creates a new web search client
func NewWebSearchClient(token string) *WebSearchClient {
	if token == "" {
		token = defaultToken
	}
	return &WebSearchClient{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Search performs a web search using BigModel API
func (c *WebSearchClient) Search(ctx context.Context, query string, searchIntent bool) (*WebSearchResponse, error) {
	reqBody := WebSearchRequest{
		SearchQuery:  query,
		SearchEngine: "search_std",
		SearchIntent: searchIntent,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", webSearchURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp WebSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResp, nil
}

// NewWebFetchClient creates a new web fetch client
func NewWebFetchClient() *WebFetchClient {
	return &WebFetchClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchWebPage fetches and extracts content from a web page
func (c *WebFetchClient) FetchWebPage(ctx context.Context, targetURL string) (*WebPageContent, error) {
	// Validate URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow HTTP and HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; EZ-Web-Search-MCP/1.0)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract content
	content := &WebPageContent{
		URL:         targetURL,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Headers:     make(map[string]string),
	}

	// Extract basic headers
	for key, values := range resp.Header {
		if len(values) > 0 {
			content.Headers[key] = values[0]
		}
	}

	// Extract title
	content.Title = doc.Find("title").First().Text()
	content.Title = strings.TrimSpace(content.Title)

	// Extract meta description
	content.Description = doc.Find("meta[name='description']").AttrOr("content", "")
	if content.Description == "" {
		content.Description = doc.Find("meta[property='og:description']").AttrOr("content", "")
	}

	// Extract meta keywords
	content.Keywords = doc.Find("meta[name='keywords']").AttrOr("content", "")

	// Extract author
	content.Author = doc.Find("meta[name='author']").AttrOr("content", "")
	if content.Author == "" {
		content.Author = doc.Find("meta[property='article:author']").AttrOr("content", "")
	}

	// Extract language
	content.Language = doc.Find("html").AttrOr("lang", "")
	if content.Language == "" {
		content.Language = doc.Find("meta[http-equiv='content-language']").AttrOr("content", "")
	}

	// Extract main content (try multiple selectors)
	var textContent strings.Builder

	// Try common content selectors
	contentSelectors := []string{
		"article", "main", ".content", "#content", ".post-content",
		".entry-content", ".article-content", ".post-body", "p",
	}

	for _, selector := range contentSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if len(text) > 50 { // Only include substantial text blocks
				textContent.WriteString(text)
				textContent.WriteString("\n\n")
			}
		})
		if textContent.Len() > 500 { // Stop if we have enough content
			break
		}
	}

	content.Content = strings.TrimSpace(textContent.String())

	// If no content found, extract all text
	if content.Content == "" {
		content.Content = strings.TrimSpace(doc.Find("body").Text())
	}

	// Clean up content (remove excessive whitespace)
	re := regexp.MustCompile(`\s+`)
	content.Content = re.ReplaceAllString(content.Content, " ")

	// Limit content length
	if len(content.Content) > 5000 {
		content.Content = content.Content[:5000] + "..."
	}

	// Extract links
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && href != "" {
			// Convert relative URLs to absolute
			if absoluteURL, err := parsedURL.Parse(href); err == nil {
				content.Links = append(content.Links, absoluteURL.String())
			}
		}
	})

	// Extract images
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			// Convert relative URLs to absolute
			if absoluteURL, err := parsedURL.Parse(src); err == nil {
				content.Images = append(content.Images, absoluteURL.String())
			}
		}
	})

	// Limit arrays to prevent excessive data
	if len(content.Links) > 50 {
		content.Links = content.Links[:50]
	}
	if len(content.Images) > 20 {
		content.Images = content.Images[:20]
	}

	return content, nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"EZ Web Search & Fetch MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Initialize clients
	searchClient := NewWebSearchClient("")
	fetchClient := NewWebFetchClient()

	// Add web search tool
	webSearchTool := mcp.NewTool("web_search",
		mcp.WithDescription("Search the web using BigModel Web Search API"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query to execute"),
		),
		mcp.WithBoolean("search_intent",
			mcp.Description("Whether to enable search intent analysis (default: false)"),
		),
	)

	// Add web search handler
	s.AddTool(webSearchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract query parameter
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid query parameter: %v", err)), nil
		}

		// Extract search_intent parameter (optional, defaults to false)
		searchIntent := false
		if intentVal, exists := request.GetArguments()["search_intent"]; exists {
			if boolVal, ok := intentVal.(bool); ok {
				searchIntent = boolVal
			}
		}

		// Perform the search
		searchResp, err := searchClient.Search(ctx, query, searchIntent)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		// Format the response
		var resultText string
		resultText += fmt.Sprintf("Search Results for: %s\n", query)
		resultText += fmt.Sprintf("Request ID: %s\n\n", searchResp.RequestID)

		// Add search intent information if available
		if len(searchResp.SearchIntent) > 0 {
			resultText += "Search Intent Analysis:\n"
			for _, intent := range searchResp.SearchIntent {
				resultText += fmt.Sprintf("- Query: %s\n", intent.Query)
				resultText += fmt.Sprintf("  Intent: %s\n", intent.Intent)
				resultText += fmt.Sprintf("  Keywords: %s\n\n", intent.Keywords)
			}
		}

		// Add search results
		if len(searchResp.SearchResult) > 0 {
			resultText += "Search Results:\n"
			for i, result := range searchResp.SearchResult {
				resultText += fmt.Sprintf("%d. %s\n", i+1, result.Title)
				resultText += fmt.Sprintf("   URL: %s\n", result.Link)
				if result.Content != "" {
					resultText += fmt.Sprintf("   Summary: %s\n", result.Content)
				}
				if result.PublishDate != "" {
					resultText += fmt.Sprintf("   Published: %s\n", result.PublishDate)
				}
				if result.Refer != "" {
					resultText += fmt.Sprintf("   Source: %s\n", result.Refer)
				}
				resultText += "\n"
			}
		} else {
			resultText += "No search results found.\n"
		}

		return mcp.NewToolResultText(resultText), nil
	})

	// Add web fetch tool
	webFetchTool := mcp.NewTool("web_fetch",
		mcp.WithDescription("Fetch and extract content from a web page"),
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

	// Add web fetch handler
	s.AddTool(webFetchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		content, err := fetchClient.FetchWebPage(ctx, targetURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch web page: %v", err)), nil
		}

		// Format the response
		var resultText string
		resultText += fmt.Sprintf("Web Page Content for: %s\n", content.URL)
		resultText += fmt.Sprintf("Status Code: %d\n", content.StatusCode)
		resultText += fmt.Sprintf("Content Type: %s\n\n", content.ContentType)

		if content.Title != "" {
			resultText += fmt.Sprintf("Title: %s\n\n", content.Title)
		}

		if content.Description != "" {
			resultText += fmt.Sprintf("Description: %s\n\n", content.Description)
		}

		if content.Author != "" {
			resultText += fmt.Sprintf("Author: %s\n", content.Author)
		}

		if content.Language != "" {
			resultText += fmt.Sprintf("Language: %s\n", content.Language)
		}

		if content.Keywords != "" {
			resultText += fmt.Sprintf("Keywords: %s\n", content.Keywords)
		}

		if content.Author != "" || content.Language != "" || content.Keywords != "" {
			resultText += "\n"
		}

		if content.Content != "" {
			resultText += fmt.Sprintf("Content:\n%s\n\n", content.Content)
		}

		if includeLinks && len(content.Links) > 0 {
			resultText += fmt.Sprintf("Links (%d found):\n", len(content.Links))
			for i, link := range content.Links {
				if i >= 10 { // Limit to first 10 links in output
					resultText += fmt.Sprintf("... and %d more links\n", len(content.Links)-10)
					break
				}
				resultText += fmt.Sprintf("- %s\n", link)
			}
			resultText += "\n"
		}

		if includeImages && len(content.Images) > 0 {
			resultText += fmt.Sprintf("Images (%d found):\n", len(content.Images))
			for i, image := range content.Images {
				if i >= 5 { // Limit to first 5 images in output
					resultText += fmt.Sprintf("... and %d more images\n", len(content.Images)-5)
					break
				}
				resultText += fmt.Sprintf("- %s\n", image)
			}
			resultText += "\n"
		}

		return mcp.NewToolResultText(resultText), nil
	})

	// Add a simple ping tool for testing
	pingTool := mcp.NewTool("ping",
		mcp.WithDescription("Simple ping tool to test server connectivity"),
	)

	s.AddTool(pingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("pong"), nil
	})

	// Start the stdio server
	log.Println("Starting EZ Web Search & Fetch MCP Server...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

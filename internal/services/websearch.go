package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ez-web-search/internal/config"
	"ez-web-search/internal/utils"
	"ez-web-search/pkg/types"
)

// WebSearchService handles web search operations
type WebSearchService struct {
	config     *config.Config
	httpClient *http.Client
	antiBot    *utils.AntiBotManager
}

// NewWebSearchService creates a new web search service
func NewWebSearchService(cfg *config.Config) *WebSearchService {
	return &WebSearchService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.BigModel.Timeout,
		},
		antiBot: utils.NewAntiBotManager(cfg.UserAgent.Pool),
	}
}

// Search performs a web search using BigModel API
func (s *WebSearchService) Search(ctx context.Context, opts types.WebSearchOptions) (*types.WebSearchResponse, error) {
	// Use provided search engine or fall back to config default
	searchEngine := opts.SearchEngine
	if searchEngine == "" {
		searchEngine = s.config.BigModel.SearchEngine
	}

	reqBody := types.WebSearchRequest{
		SearchQuery:  opts.Query,
		SearchEngine: searchEngine,
		SearchIntent: opts.SearchIntent,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BigModel.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication and content type
	req.Header.Set("Authorization", "Bearer "+s.config.BigModel.Token)
	req.Header.Set("Content-Type", "application/json")

	// Set anti-bot headers
	userAgent := s.antiBot.GetRandomUserAgent()
	s.antiBot.SetRealisticHeaders(req, userAgent)

	// Apply random delay if configured
	if s.config.WebFetch.UserAgentRotate && s.antiBot.ShouldDelay() {
		delay := s.antiBot.GetRandomDelay(s.config.WebFetch.DelayMin, s.config.WebFetch.DelayMax)
		time.Sleep(delay)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check for rate limiting or blocking
	if s.antiBot.IsRateLimited(resp) {
		return nil, fmt.Errorf("request was rate limited, status: %d", resp.StatusCode)
	}

	if s.antiBot.IsBlocked(resp) {
		return nil, fmt.Errorf("request was blocked, status: %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp types.WebSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResp, nil
}

// FormatSearchResponse formats the search response for display
func (s *WebSearchService) FormatSearchResponse(resp *types.WebSearchResponse, query string, searchEngine string) string {
	var resultText string
	resultText += fmt.Sprintf("Search Results for: %s\n", query)
	resultText += fmt.Sprintf("Search Engine: %s\n", searchEngine)
	resultText += fmt.Sprintf("Request ID: %s\n\n", resp.RequestID)

	// Add search intent information if available
	if len(resp.SearchIntent) > 0 {
		resultText += "Search Intent Analysis:\n"
		for _, intent := range resp.SearchIntent {
			resultText += fmt.Sprintf("- Query: %s\n", intent.Query)
			resultText += fmt.Sprintf("  Intent: %s\n", intent.Intent)
			resultText += fmt.Sprintf("  Keywords: %s\n\n", intent.Keywords)
		}
	}

	// Add search results
	if len(resp.SearchResult) > 0 {
		resultText += "Search Results:\n"
		for i, result := range resp.SearchResult {
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

	return resultText
}

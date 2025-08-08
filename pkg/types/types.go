package types

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

// WebFetchOptions represents options for web fetching
type WebFetchOptions struct {
	URL           string
	IncludeLinks  bool
	IncludeImages bool
	UserAgent     string
}

// WebSearchOptions represents options for web searching
type WebSearchOptions struct {
	Query        string
	SearchIntent bool
}

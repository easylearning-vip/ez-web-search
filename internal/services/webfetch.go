package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"ez-web-search/internal/config"
	"ez-web-search/internal/utils"
	"ez-web-search/pkg/types"
)

// WebFetchService handles web page fetching and content extraction
type WebFetchService struct {
	config     *config.Config
	httpClient *http.Client
	antiBot    *utils.AntiBotManager
}

// NewWebFetchService creates a new web fetch service
func NewWebFetchService(cfg *config.Config) *WebFetchService {
	return &WebFetchService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.WebFetch.Timeout,
		},
		antiBot: utils.NewAntiBotManager(cfg.UserAgent.Pool),
	}
}

// FetchWebPage fetches and extracts content from a web page with anti-bot measures
func (s *WebFetchService) FetchWebPage(ctx context.Context, opts types.WebFetchOptions) (*types.WebPageContent, error) {
	// Validate URL
	parsedURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow HTTP and HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	// Apply random delay to avoid detection
	if s.config.WebFetch.UserAgentRotate && s.antiBot.ShouldDelay() {
		delay := s.antiBot.GetRandomDelay(s.config.WebFetch.DelayMin, s.config.WebFetch.DelayMax)
		time.Sleep(delay)
	}

	// Create request with random timeout variance
	timeout := s.antiBot.GetRandomTimeout(s.config.WebFetch.Timeout)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxWithTimeout, "GET", opts.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Use provided user agent or get a random one
	userAgent := opts.UserAgent
	if userAgent == "" {
		userAgent = s.antiBot.GetRandomUserAgent()
	}

	// Set realistic headers to avoid bot detection
	s.antiBot.SetRealisticHeaders(req, userAgent)

	// Perform request with retry logic
	var resp *http.Response
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = s.httpClient.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to fetch page after %d attempts: %w", maxRetries, err)
			}
			// Wait before retry
			retryDelay := s.antiBot.GetRetryDelay(resp, attempt)
			time.Sleep(retryDelay)
			continue
		}

		// Check for rate limiting
		if s.antiBot.IsRateLimited(resp) {
			resp.Body.Close()
			if attempt == maxRetries {
				return nil, fmt.Errorf("request was rate limited after %d attempts", maxRetries)
			}
			retryDelay := s.antiBot.GetRetryDelay(resp, attempt)
			time.Sleep(retryDelay)
			continue
		}

		// Check for blocking
		if s.antiBot.IsBlocked(resp) {
			resp.Body.Close()
			return nil, fmt.Errorf("request was blocked, status: %d", resp.StatusCode)
		}

		// Success, break out of retry loop
		break
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if response is gzip compressed and decompress if needed
	var reader io.Reader = bytes.NewReader(body)
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") || len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b {
		gzipReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		decompressed, err := io.ReadAll(gzipReader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response: %w", err)
		}
		reader = bytes.NewReader(decompressed)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract content
	content := &types.WebPageContent{
		URL:         opts.URL,
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

	// Extract metadata
	s.extractMetadata(doc, content)

	// Extract main content
	s.extractContent(doc, content)

	// Extract links if requested
	if opts.IncludeLinks {
		s.extractLinks(doc, content, parsedURL)
	}

	// Extract images if requested
	if opts.IncludeImages {
		s.extractImages(doc, content, parsedURL)
	}

	return content, nil
}

// extractMetadata extracts metadata from the HTML document
func (s *WebFetchService) extractMetadata(doc *goquery.Document, content *types.WebPageContent) {
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
}

// extractContent extracts main content from the HTML document
func (s *WebFetchService) extractContent(doc *goquery.Document, content *types.WebPageContent) {
	var textContent strings.Builder

	// Try common content selectors in order of preference
	contentSelectors := []string{
		"article", "main", ".content", "#content", ".post-content",
		".entry-content", ".article-content", ".post-body", ".article-body",
		"[role='main']", ".main-content", "#main-content",
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

	// If no content found, extract paragraphs
	if textContent.Len() < 100 {
		doc.Find("p").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if len(text) > 30 {
				textContent.WriteString(text)
				textContent.WriteString("\n\n")
			}
		})
	}

	// If still no content, extract all text from body
	if textContent.Len() < 50 {
		bodyText := strings.TrimSpace(doc.Find("body").Text())
		textContent.WriteString(bodyText)
	}

	content.Content = strings.TrimSpace(textContent.String())

	// Clean up content (remove excessive whitespace)
	re := regexp.MustCompile(`\s+`)
	content.Content = re.ReplaceAllString(content.Content, " ")

	// Limit content length
	if len(content.Content) > s.config.WebFetch.MaxContentSize {
		content.Content = content.Content[:s.config.WebFetch.MaxContentSize] + "..."
	}
}

// extractLinks extracts links from the HTML document
func (s *WebFetchService) extractLinks(doc *goquery.Document, content *types.WebPageContent, baseURL *url.URL) {
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && href != "" && href != "#" {
			// Convert relative URLs to absolute
			if absoluteURL, err := baseURL.Parse(href); err == nil {
				content.Links = append(content.Links, absoluteURL.String())
			}
		}
	})

	// Limit links to prevent excessive data
	if len(content.Links) > s.config.WebFetch.MaxLinks {
		content.Links = content.Links[:s.config.WebFetch.MaxLinks]
	}
}

// extractImages extracts images from the HTML document
func (s *WebFetchService) extractImages(doc *goquery.Document, content *types.WebPageContent, baseURL *url.URL) {
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			// Convert relative URLs to absolute
			if absoluteURL, err := baseURL.Parse(src); err == nil {
				content.Images = append(content.Images, absoluteURL.String())
			}
		}
	})

	// Limit images to prevent excessive data
	if len(content.Images) > s.config.WebFetch.MaxImages {
		content.Images = content.Images[:s.config.WebFetch.MaxImages]
	}
}

// FormatWebPageContent formats the web page content for display
func (s *WebFetchService) FormatWebPageContent(content *types.WebPageContent, includeLinks, includeImages bool) string {
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

	return resultText
}

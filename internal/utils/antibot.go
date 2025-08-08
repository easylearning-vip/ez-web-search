package utils

import (
	"math/rand"
	"net/http"
	"time"
)

// AntiBotManager handles anti-bot detection bypass mechanisms
type AntiBotManager struct {
	userAgents []string
	rand       *rand.Rand
}

// NewAntiBotManager creates a new anti-bot manager
func NewAntiBotManager(userAgents []string) *AntiBotManager {
	return &AntiBotManager{
		userAgents: userAgents,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetRandomUserAgent returns a random user agent from the pool
func (a *AntiBotManager) GetRandomUserAgent() string {
	if len(a.userAgents) == 0 {
		return "Mozilla/5.0 (compatible; EZ-Web-Search-MCP/1.0)"
	}
	return a.userAgents[a.rand.Intn(len(a.userAgents))]
}

// SetRealisticHeaders sets realistic browser headers to avoid detection
func (a *AntiBotManager) SetRealisticHeaders(req *http.Request, userAgent string) {
	// Set user agent
	req.Header.Set("User-Agent", userAgent)
	
	// Set realistic browser headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	
	// Vary some headers to look more natural
	secFetchHeaders := [][]string{
		{"Sec-Fetch-Dest", "document"},
		{"Sec-Fetch-Mode", "navigate"},
		{"Sec-Fetch-Site", "none"},
		{"Sec-Fetch-User", "?1"},
	}
	
	for _, header := range secFetchHeaders {
		req.Header.Set(header[0], header[1])
	}
	
	// Add cache control occasionally
	if a.rand.Float32() < 0.3 {
		req.Header.Set("Cache-Control", "max-age=0")
	}
	
	// Add referer occasionally to simulate navigation
	if a.rand.Float32() < 0.2 {
		referers := []string{
			"https://www.google.com/",
			"https://www.bing.com/",
			"https://duckduckgo.com/",
		}
		req.Header.Set("Referer", referers[a.rand.Intn(len(referers))])
	}
}

// GetRandomDelay returns a random delay between min and max duration
func (a *AntiBotManager) GetRandomDelay(min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}
	
	minNano := min.Nanoseconds()
	maxNano := max.Nanoseconds()
	randomNano := minNano + a.rand.Int63n(maxNano-minNano)
	
	return time.Duration(randomNano)
}

// ShouldDelay determines if a delay should be applied (with some randomness)
func (a *AntiBotManager) ShouldDelay() bool {
	// 70% chance to apply delay
	return a.rand.Float32() < 0.7
}

// GetRandomTimeout returns a random timeout with some variance
func (a *AntiBotManager) GetRandomTimeout(baseTimeout time.Duration) time.Duration {
	// Add ±20% variance to timeout
	variance := float64(baseTimeout) * 0.2
	adjustment := (a.rand.Float64() - 0.5) * 2 * variance
	return time.Duration(float64(baseTimeout) + adjustment)
}

// IsRateLimited checks if the response indicates rate limiting
func (a *AntiBotManager) IsRateLimited(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	
	// Common rate limiting status codes
	rateLimitCodes := []int{429, 503, 509}
	for _, code := range rateLimitCodes {
		if resp.StatusCode == code {
			return true
		}
	}
	
	// Check for rate limiting headers
	rateLimitHeaders := []string{
		"X-RateLimit-Remaining",
		"X-Rate-Limit-Remaining",
		"RateLimit-Remaining",
		"Retry-After",
	}
	
	for _, header := range rateLimitHeaders {
		if resp.Header.Get(header) != "" {
			return true
		}
	}
	
	return false
}

// GetRetryDelay calculates appropriate retry delay for rate limited requests
func (a *AntiBotManager) GetRetryDelay(resp *http.Response, attempt int) time.Duration {
	// Check for Retry-After header
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		if duration, err := time.ParseDuration(retryAfter + "s"); err == nil {
			return duration
		}
	}
	
	// Exponential backoff with jitter
	baseDelay := time.Duration(attempt) * time.Second
	maxDelay := 30 * time.Second
	
	if baseDelay > maxDelay {
		baseDelay = maxDelay
	}
	
	// Add jitter (±25%)
	jitter := float64(baseDelay) * 0.25
	adjustment := (a.rand.Float64() - 0.5) * 2 * jitter
	
	return time.Duration(float64(baseDelay) + adjustment)
}

// IsBlocked checks if the response indicates the request was blocked
func (a *AntiBotManager) IsBlocked(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	
	// Common blocking status codes
	blockingCodes := []int{403, 406, 418, 451}
	for _, code := range blockingCodes {
		if resp.StatusCode == code {
			return true
		}
	}
	
	// Check for common blocking indicators in headers
	blockingHeaders := map[string][]string{
		"Server": {"cloudflare", "nginx"},
		"CF-Ray": {""},  // Cloudflare
		"X-Sucuri-ID": {""},  // Sucuri WAF
	}
	
	for header, values := range blockingHeaders {
		headerValue := resp.Header.Get(header)
		if headerValue != "" {
			for _, value := range values {
				if value == "" || headerValue == value {
					// Additional check for Cloudflare challenge
					if header == "Server" && headerValue == "cloudflare" {
						if resp.StatusCode == 403 || resp.StatusCode == 503 {
							return true
						}
					}
				}
			}
		}
	}
	
	return false
}

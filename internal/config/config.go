package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	BigModel  BigModelConfig
	WebFetch  WebFetchConfig
	UserAgent UserAgentConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Name    string
	Version string
}

// BigModelConfig holds BigModel API configuration
type BigModelConfig struct {
	Token        string
	BaseURL      string
	Timeout      time.Duration
	SearchEngine string
}

// WebFetchConfig holds web fetching configuration
type WebFetchConfig struct {
	Timeout         time.Duration
	MaxContentSize  int
	MaxLinks        int
	MaxImages       int
	UserAgentRotate bool
	DelayMin        time.Duration
	DelayMax        time.Duration
}

// UserAgentConfig holds user agent rotation configuration
type UserAgentConfig struct {
	Pool []string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Name:    getEnv("SERVER_NAME", "EZ Web Search & Fetch MCP Server"),
			Version: getEnv("SERVER_VERSION", "1.0.0"),
		},
		BigModel: BigModelConfig{
			Token:        getEnv("BIGMODEL_TOKEN", ""),
			BaseURL:      getEnv("BIGMODEL_BASE_URL", "https://open.bigmodel.cn/api/paas/v4/web_search"),
			Timeout:      getDurationEnv("BIGMODEL_TIMEOUT", 30*time.Second),
			SearchEngine: getEnv("BIGMODEL_SEARCH_ENGINE", "search_std"),
		},
		WebFetch: WebFetchConfig{
			Timeout:         getDurationEnv("WEBFETCH_TIMEOUT", 30*time.Second),
			MaxContentSize:  getIntEnv("WEBFETCH_MAX_CONTENT_SIZE", 5000),
			MaxLinks:        getIntEnv("WEBFETCH_MAX_LINKS", 50),
			MaxImages:       getIntEnv("WEBFETCH_MAX_IMAGES", 20),
			UserAgentRotate: getBoolEnv("WEBFETCH_USER_AGENT_ROTATE", true),
			DelayMin:        getDurationEnv("WEBFETCH_DELAY_MIN", 1*time.Second),
			DelayMax:        getDurationEnv("WEBFETCH_DELAY_MAX", 3*time.Second),
		},
		UserAgent: UserAgentConfig{
			Pool: getDefaultUserAgents(),
		},
	}
}

// Validate validates the configuration and returns an error if invalid
func (c *Config) Validate() error {
	if c.BigModel.Token == "" {
		return fmt.Errorf("BIGMODEL_TOKEN is required and must be set")
	}
	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable with a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getBoolEnv gets a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable with a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getDefaultUserAgents returns a pool of realistic user agents for rotation
func getDefaultUserAgents() []string {
	return []string{
		// Chrome on Windows
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",

		// Chrome on macOS
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",

		// Firefox on Windows
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",

		// Firefox on macOS
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",

		// Safari on macOS
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",

		// Edge on Windows
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",

		// Chrome on Linux
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",

		// Mobile Chrome
		"Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	}
}

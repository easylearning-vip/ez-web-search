# EZ Web Search & Fetch MCP Server Makefile

.PHONY: build test clean run dev install-deps lint format help

# Variables
BINARY_NAME=ez-web-search-v2
BINARY_PATH=./$(BINARY_NAME)
CMD_PATH=./cmd/server
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
help: ## Show this help message
	@echo "EZ Web Search & Fetch MCP Server"
	@echo "================================"
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the server binary
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_PATH)"

build-release: ## Build optimized release binary
	@echo "Building release version..."
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Release build complete: $(BINARY_PATH)"

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean

run: build ## Build and run the server
	@echo "Starting server..."
	$(BINARY_PATH)

dev: ## Run in development mode with environment variables
	@echo "Starting development server..."
	@if [ -f .env ]; then \
		export $$(cat .env | xargs) && $(BINARY_PATH); \
	else \
		echo "No .env file found. Copy .env.example to .env and configure it."; \
		$(BINARY_PATH); \
	fi

install-deps: ## Install/update dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

format: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w $(GO_FILES); \
	else \
		echo "goimports not installed. Install it with:"; \
		echo "go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

# Testing targets
test-ping: build ## Test ping tool
	@echo "Testing ping tool..."
	npx @modelcontextprotocol/inspector --cli $(BINARY_PATH) --method tools/call --tool-name ping

test-search: build ## Test web search tool
	@echo "Testing web search tool..."
	npx @modelcontextprotocol/inspector --cli $(BINARY_PATH) --method tools/call --tool-name web_search --tool-arg query="Go programming tutorial"

test-fetch: build ## Test web fetch tool
	@echo "Testing web fetch tool..."
	npx @modelcontextprotocol/inspector --cli $(BINARY_PATH) --method tools/call --tool-name web_fetch --tool-arg url="http://httpbin.org/html"

test-all-tools: build ## Test all tools
	@echo "Testing all tools..."
	@echo "\n1. Testing tools list..."
	npx @modelcontextprotocol/inspector --cli $(BINARY_PATH) --method tools/list
	@echo "\n2. Testing ping..."
	npx @modelcontextprotocol/inspector --cli $(BINARY_PATH) --method tools/call --tool-name ping
	@echo "\n3. Testing web search..."
	npx @modelcontextprotocol/inspector --cli $(BINARY_PATH) --method tools/call --tool-name web_search --tool-arg query="MCP testing"
	@echo "\n4. Testing web fetch..."
	npx @modelcontextprotocol/inspector --cli $(BINARY_PATH) --method tools/call --tool-name web_fetch --tool-arg url="http://httpbin.org/html"

inspector-ui: build ## Start MCP Inspector UI
	@echo "Starting MCP Inspector UI..."
	npx @modelcontextprotocol/inspector $(BINARY_PATH)

# Configuration targets
config-example: ## Show configuration example
	@echo "Configuration example:"
	@cat .env.example

# Development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Project structure
show-structure: ## Show project structure
	@echo "Project structure:"
	@tree -I 'node_modules|.git|*.sum' || find . -type f -name "*.go" | head -20

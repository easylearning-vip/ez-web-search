# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an enterprise-grade Web Search & Fetch MCP (Model Context Protocol) server implemented in Go. It provides AI applications with comprehensive web interaction capabilities including web search via BigModel API and intelligent web content extraction with anti-bot protection.

## Architecture

### Core Structure
- **Entry Point**: `cmd/server/main.go` - Server initialization and configuration
- **MCP Handlers**: `internal/handlers/mcp.go` - Tool registration and request handling
- **Business Logic**: 
  - `internal/services/websearch.go` - BigModel API integration
  - `internal/services/webfetch.go` - Web content extraction
- **Configuration**: `internal/config/config.go` - Environment-based config management
- **Utilities**: `internal/utils/antibot.go` - Anti-bot protection mechanisms
- **Types**: `pkg/types/types.go` - Shared type definitions

### MCP Protocol Implementation
- Uses `mark3labs/mcp-go v0.37.0` library
- Implements MCP 2024-11-05 protocol specification
- Supports stdio transport for client integration
- Three main tools: `ping`, `ez_web_search`, `ez_web_fetch`

## Development Commands

### Building and Running
```bash
# Build the server
make build

# Build optimized release version
make build-release

# Run the server
make run

# Run in development mode (loads .env file)
make dev
```

### Testing
```bash
# Run Go tests
make test

# Run tests with coverage
make test-coverage

# Test individual MCP tools
make test-ping          # Test ping tool
make test-search        # Test web search tool
make test-fetch         # Test web fetch tool
make test-all-tools     # Test all tools

# Start interactive MCP Inspector UI
make inspector-ui
```

### Code Quality
```bash
# Format code
make format

# Run linter (requires golangci-lint)
make lint

# Install development tools
make install-tools
```

### MCP Protocol Testing
```bash
# Official MCP Inspector CLI tests
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/list
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name ping
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name ez_web_search --tool-arg query="test query"
```

## Configuration

### Environment Variables
- `BIGMODEL_TOKEN` - Required API token for web search
- `BIGMODEL_BASE_URL` - BigModel API base URL (default: https://open.bigmodel.cn/api/paas/v4/web_search)
- `BIGMODEL_TIMEOUT` - API timeout (default: 30s)
- `WEBFETCH_TIMEOUT` - Web fetch timeout (default: 30s)
- `WEBFETCH_MAX_CONTENT_SIZE` - Max content size to fetch (default: 5000)
- `WEBFETCH_USER_AGENT_ROTATE` - Enable user agent rotation (default: true)

### Command Line Options
- `-token string` - Override BigModel API token
- `-help` - Show usage information

## Key Implementation Details

### Web Search Tool (`ez_web_search`)
- Integrates with BigModel Web Search API
- Supports multiple search engines: search_std, search_pro, search_pro_sogou, search_pro_quark
- Optional search intent analysis and keyword extraction
- Returns structured results with titles, URLs, summaries, publication dates

### Web Fetch Tool (`ez_web_fetch`)
- Intelligent content extraction using goquery
- Extracts metadata: title, description, author, keywords, language
- Optional link and image extraction with relative URL resolution
- Anti-bot protection with user agent rotation and request delays
- Security validation (HTTP/HTTPS only)
- Content length limiting (5000 chars max)

### Anti-Bot Protection
- User agent rotation with 12+ real browser user agents
- Complete browser header spoofing
- Configurable random delays (1-3 seconds)
- Exponential backoff retry logic
- Rate limiting detection (429, 503 status codes)

## Testing Strategy

The project uses multiple testing approaches:

1. **Official MCP Inspector** - Primary testing for protocol compliance
2. **Go unit tests** - Business logic testing
3. **Integration tests** - End-to-end functionality verification
4. **Manual testing** - Development and debugging

Test scripts are located in project root and use the official MCP Inspector tool for comprehensive validation.

## Dependencies

- **Go 1.23** - Core language and runtime
- **mark3labs/mcp-go v0.37.0** - MCP protocol implementation
- **PuerkitoBio/goquery v1.10.3** - HTML parsing and manipulation
- **BigModel Web Search API** - External search service

## Build System

Uses a comprehensive Makefile with targets for:
- Building (debug/release)
- Testing (unit, integration, MCP protocol)
- Code quality (linting, formatting)
- Development tools installation
- Release management

The project supports cross-platform builds and has automated release scripts.
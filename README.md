# EZ Web Search & Fetch MCP Server

A complete, enterprise-grade Web Search and Fetch MCP (Model Context Protocol) server implemented in Go, providing comprehensive web interaction capabilities for AI applications with advanced anti-bot protection.

## 🚀 Features

- **Web Search Tool**: Search the web using BigModel's Web Search API
- **Web Fetch Tool**: Fetch and extract content from any web page with anti-bot protection
- **Search Intent Analysis**: Optional search intent analysis and keyword extraction
- **Content Extraction**: Intelligent extraction of titles, descriptions, text content, links, and images
- **Anti-Bot Protection**: Advanced mechanisms to bypass detection and rate limiting
- **Enterprise Architecture**: Modular, scalable design following Go best practices
- **Environment Configuration**: Secure token management via environment variables
- **MCP Protocol Compliance**: Built using the official mark3labs/mcp-go library
- **Docker Support**: Containerized deployment ready
- **Comprehensive Testing**: Multiple testing approaches and tools

## Prerequisites

- Go 1.21 or later
- BigModel API token (default test token included)

## 🏗️ Enterprise Architecture

### Project Structure
```
ez-web-search/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # MCP request handlers
│   ├── services/        # Business logic services
│   └── utils/           # Utility functions (anti-bot, etc.)
├── pkg/
│   ├── client/          # External API clients
│   └── types/           # Shared type definitions
├── scripts/             # Build and deployment scripts
├── docs/               # Documentation
├── Makefile           # Build automation
├── .gitignore         # Git ignore rules
└── .env.example       # Configuration template
```

### Technology Stack

- **Go 1.23**: Modern Go with latest features and performance
- **mark3labs/mcp-go v0.37.0**: Official MCP Go library
- **PuerkitoBio/goquery v1.10.3**: jQuery-like HTML parsing and manipulation
- **BigModel Web Search API**: Professional web search service
- **Anti-Bot Protection**: User agent rotation, request delays, header spoofing
- **Environment Configuration**: Secure configuration via environment variables
- **Git Integration**: Version control ready with proper .gitignore
- **Standard Go libraries**: net/http, context, encoding/json, regexp, etc.

## 🚀 Quick Start

### Using Make (Recommended)

```bash
# Build the server
make build

# Run with default configuration
make run

# Run all tests
make test-all-tools

# Start Inspector UI
make inspector-ui
```

### Manual Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd ez-web-search
```

2. Install dependencies:
```bash
make install-deps
# or manually: go mod download && go mod tidy
```

3. Configure environment (optional):
```bash
cp .env.example .env
# Edit .env with your BigModel API token
```

4. Build and run:
```bash
make build
./ez-web-search-v2
```

### Environment Configuration

```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your configuration
# BIGMODEL_TOKEN="your_actual_bigmodel_api_token"
# WEBFETCH_USER_AGENT_ROTATE=true
# WEBFETCH_DELAY_MIN="1s"
# WEBFETCH_DELAY_MAX="3s"

# Run with environment variables
make dev
```

## Usage

### Running the Server

The server runs in stdio mode, which is the standard way to run MCP servers:

```bash
./ez-web-search
```

### Available Tools

#### 1. web_search
Searches the web using BigModel's Web Search API.

**Parameters:**
- `query` (required, string): The search query to execute
- `search_intent` (optional, boolean): Whether to enable search intent analysis (default: false)

**Example usage in MCP client:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "web_search",
    "arguments": {
      "query": "Go programming language tutorial",
      "search_intent": true
    }
  }
}
```

#### 2. web_fetch
Fetches and extracts content from any web page.

**Parameters:**
- `url` (required, string): The URL of the web page to fetch
- `include_links` (optional, boolean): Whether to include extracted links (default: false)
- `include_images` (optional, boolean): Whether to include extracted images (default: false)

**Features:**
- Extracts title, description, author, language, keywords
- Intelligent content extraction from articles and main content areas
- Converts relative URLs to absolute URLs
- Handles various HTML structures and meta tags
- Content length limiting and cleanup

**Example usage:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "web_fetch",
    "arguments": {
      "url": "https://example.com/article",
      "include_links": true,
      "include_images": false
    }
  }
}
```

#### 3. ping
Simple connectivity test tool.

**Parameters:** None

**Example usage:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "ping",
    "arguments": {}
  }
}
```

## Configuration

### API Token

The server uses a default test token. To use your own token, modify the `defaultToken` constant in `main.go` or set it via environment variable (future enhancement).

### Search Engine

The server uses `search_std` as the default search engine. This can be modified in the `Search` method of `WebSearchClient`.

## API Reference

This server integrates with BigModel's Web Search API:
- **Endpoint**: `https://open.bigmodel.cn/api/paas/v4/web_search`
- **Documentation**: https://docs.bigmodel.cn/api-reference/工具-api/网络搜索

## Response Format

The web search tool returns formatted text containing:
- Search query and request ID
- Search intent analysis (if enabled)
- Search results with titles, URLs, summaries, and metadata

## Testing

### 🔧 Official MCP Inspector Testing (Recommended)

The **MCP Inspector** is the official testing tool from the Model Context Protocol team. It provides both CLI and UI modes for comprehensive testing.

#### CLI Mode Testing

Quick command-line testing for automation and scripting:

```bash
# List available tools
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/list

# Test ping tool
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name ping

# Test web search
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="Go programming tutorial"

# Test web fetch
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_fetch --tool-arg url="https://example.com"

# Test web fetch with links and images
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_fetch --tool-arg url="https://example.com" --tool-arg include_links=true --tool-arg include_images=true

# Test with search intent analysis
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="MCP testing" --tool-arg search_intent=true
```

#### UI Mode Testing

Interactive web-based testing interface:

```bash
# Start the Inspector UI (opens browser automatically)
npx @modelcontextprotocol/inspector ./ez-web-search
```

This opens a web interface at `http://localhost:6274` where you can:
- Visually inspect all available tools
- Test tools with form-based input
- View real-time responses and error messages
- Monitor server logs and notifications
- Export server configurations for other MCP clients

#### Automated Inspector Testing

Use our test script for guided testing:

```bash
./test_with_inspector.sh
```

### 🧪 Custom Testing Scripts

#### Automated Verification

Run the verification script to check basic functionality:

```bash
go run verify.go
```

#### Interactive Python Testing

Use the Python test script for comprehensive testing:

```bash
python3 interactive_test.py "your search query"
```

### 🔧 Manual Testing

You can test the server manually by running it and sending JSON-RPC messages:

```bash
./ez-web-search
```

Then send messages via stdin:

```json
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "ping", "arguments": {}}}
{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "web_search", "arguments": {"query": "Go programming", "search_intent": false}}}
```

### 🔗 Integration with MCP Clients

#### Using with mcphost

```bash
# Install mcphost
go install github.com/mark3labs/mcphost@latest

# Test the server
mcphost --server-command "./ez-web-search"
```

#### Using with Claude Code CLI

Add to your Claude Code CLI MCP configuration file:

**Global Configuration** (`~/.claude/mcp_settings.json`):
```json
{
  "mcpServers": {
    "ez-web-search": {
      "command": "/path/to/ez-web-search-v2",
      "env": {
        "BIGMODEL_TOKEN": "your_bigmodel_api_token",
        "WEBFETCH_USER_AGENT_ROTATE": "true",
        "PATH": "/usr/local/bin:/usr/bin:/bin"
      }
    }
  }
}
```

**Project-specific Configuration** (`.claude/mcp_settings.json`):
```json
{
  "mcpServers": {
    "ez-web-search": {
      "command": "./ez-web-search-v2",
      "env": {
        "BIGMODEL_TOKEN": "your_bigmodel_api_token"
      }
    }
  }
}
```

**Usage in Claude Code CLI**:
```bash
# Start Claude Code CLI with MCP support
claude

# Use web search in conversation
> Search for "Go web scraping best practices"

# Use web fetch in conversation
> Fetch content from https://example.com and summarize it

# Combine search and fetch
> Search for Go tutorials, then fetch content from the top 3 results
```

#### Using with Claude Desktop

Add to your Claude Desktop MCP configuration file (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "ez-web-search": {
      "command": "/path/to/ez-web-search-v2",
      "env": {
        "BIGMODEL_TOKEN": "your_bigmodel_api_token"
      }
    }
  }
}
```

#### Using with Cursor or VS Code

The Inspector UI provides export buttons to generate configuration files for various MCP clients.

### ✅ Test Results

Our server has been thoroughly tested with the official MCP Inspector:

#### Tools Discovery
```json
{
  "tools": [
    {
      "name": "ping",
      "description": "Simple ping tool to test server connectivity",
      "inputSchema": {
        "type": "object",
        "properties": {}
      }
    },
    {
      "name": "web_fetch",
      "description": "Fetch and extract content from a web page",
      "inputSchema": {
        "type": "object",
        "properties": {
          "url": {
            "description": "The URL of the web page to fetch",
            "type": "string"
          },
          "include_links": {
            "description": "Whether to include extracted links (default: false)",
            "type": "boolean"
          },
          "include_images": {
            "description": "Whether to include extracted images (default: false)",
            "type": "boolean"
          }
        },
        "required": ["url"]
      }
    },
    {
      "name": "web_search",
      "description": "Search the web using BigModel Web Search API",
      "inputSchema": {
        "type": "object",
        "properties": {
          "query": {
            "description": "The search query to execute",
            "type": "string"
          },
          "search_intent": {
            "description": "Whether to enable search intent analysis (default: false)",
            "type": "boolean"
          }
        },
        "required": ["query"]
      }
    }
  ]
}
```

#### Ping Tool Test
```json
{
  "content": [
    {
      "type": "text",
      "text": "pong"
    }
  ]
}
```

#### Web Fetch Test
```json
{
  "content": [
    {
      "type": "text",
      "text": "Web Page Content for: https://example.com\nStatus Code: 200\nContent Type: text/html\n\nTitle: Example Domain\n\nContent:\nThis domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.\n\n"
    }
  ]
}
```

#### Web Search Test
Successfully returns structured search results with:
- Search query and request ID
- Search intent analysis (when enabled)
- Multiple search results with titles, URLs, summaries
- Publication dates and source references

All tests pass with the official MCP Inspector, confirming full protocol compliance.

## 🤖 Claude Code CLI Integration

### Quick Setup for Claude Code CLI

**Automated Setup (Recommended)**:
```bash
# Run the automated setup script
./setup-claude-cli.sh
```

The script will:
- Build the MCP server
- Create Claude Code CLI configuration
- Prompt for your BigModel API token
- Test the configuration
- Provide usage examples

**Manual Setup**:

1. **Build the server**:
   ```bash
   make build
   ```

2. **Copy configuration template**:
   ```bash
   # Create Claude Code CLI MCP configuration directory
   mkdir -p ~/.claude

   # Copy and customize the configuration template
   cp claude-mcp-config.json ~/.claude/mcp_settings.json
   ```

3. **Update configuration**:
   ```bash
   # Update the binary path
   PWD_PATH=$(pwd)
   sed -i "s|/path/to/ez-web-search-v2|$PWD_PATH/ez-web-search-v2|g" ~/.claude/mcp_settings.json

   # Set your BigModel API token
   sed -i 's/your_bigmodel_api_token_here/YOUR_ACTUAL_TOKEN/g' ~/.claude/mcp_settings.json
   ```

4. **Start Claude Code CLI**:
   ```bash
   claude
   ```

### Usage Examples in Claude Code CLI

Once configured, you can use the web search and fetch tools directly in your Claude Code CLI conversations:

**Web Search Examples**:
```
> Search for "Go web scraping anti-bot techniques" and show me the latest approaches

> Find recent articles about MCP protocol implementation in Go

> Search for "PuerkitoBio goquery tutorial" with search intent analysis
```

**Web Fetch Examples**:
```
> Fetch the content from https://pkg.go.dev/github.com/PuerkitoBio/goquery and summarize the main features

> Get the documentation from https://modelcontextprotocol.io and extract the key concepts

> Fetch https://example.com with links and images included
```

**Combined Workflows**:
```
> Search for "Go HTTP client best practices", then fetch content from the top 3 results and create a comprehensive guide

> Find the latest Go web scraping libraries, fetch their documentation, and compare their features

> Search for MCP server examples, fetch the GitHub repositories, and analyze the code structure
```

### Advanced Configuration

For project-specific settings, create `.claude/mcp_settings.json` in your project directory:

```json
{
  "mcpServers": {
    "ez-web-search": {
      "command": "./ez-web-search-v2",
      "env": {
        "BIGMODEL_TOKEN": "project_specific_token",
        "WEBFETCH_MAX_CONTENT_SIZE": "10000",
        "WEBFETCH_USER_AGENT_ROTATE": "true"
      }
    }
  }
}
```

### Troubleshooting

If the MCP server doesn't work in Claude Code CLI:

1. **Check the binary path**:
   ```bash
   which ez-web-search-v2
   # or
   ls -la ./ez-web-search-v2
   ```

2. **Test the server manually**:
   ```bash
   ./ez-web-search-v2
   # Should start and wait for JSON-RPC input
   ```

3. **Verify MCP configuration**:
   ```bash
   cat ~/.claude/mcp_settings.json | jq .
   ```

4. **Check Claude Code CLI logs**:
   ```bash
   claude --debug
   ```

## 🏆 MCP Testing Best Practices

This project demonstrates MCP testing best practices:

### 1. **Multi-Level Testing Strategy**
- **Unit Testing**: Basic functionality verification (`verify.go`)
- **Integration Testing**: Full protocol testing with Inspector
- **Manual Testing**: Interactive debugging capabilities
- **Automated Testing**: CLI-based testing for CI/CD

### 2. **Official Tool Usage**
- Primary testing with **MCP Inspector** (official tool)
- CLI mode for automation and scripting
- UI mode for interactive development and debugging
- Export capabilities for client integration

### 3. **Protocol Compliance Verification**
- ✅ Proper JSON-RPC 2.0 message format
- ✅ MCP 2024-11-05 protocol version support
- ✅ Correct tool schema definitions
- ✅ Standard error handling
- ✅ Proper capability negotiation

### 4. **Development Workflow**
1. **Build** → `go build -o ez-web-search main.go`
2. **Quick Test** → `go run verify.go`
3. **Inspector CLI** → `npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/list`
4. **Inspector UI** → `npx @modelcontextprotocol/inspector ./ez-web-search`
5. **Integration** → Export config for target MCP client

### 5. **Testing Coverage**
- ✅ Server initialization and capability negotiation
- ✅ Tool discovery and schema validation (3 tools)
- ✅ Tool execution with various parameters
- ✅ Error handling and edge cases
- ✅ Real API integration (BigModel Web Search)
- ✅ Web page fetching and content extraction
- ✅ HTML parsing and metadata extraction
- ✅ Search intent analysis functionality
- ✅ URL validation and security checks

## Development

### Project Structure
```
ez-web-search/
├── main.go          # Main server implementation
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
└── README.md        # This file
```

### Key Components

1. **WebSearchClient**: Handles communication with BigModel API
2. **MCP Server**: Implements the Model Context Protocol using mark3labs/mcp-go
3. **Tool Handlers**: Process tool calls and return formatted responses

## License

This project is provided as-is for educational and testing purposes.

## Contributing

Feel free to submit issues and enhancement requests!

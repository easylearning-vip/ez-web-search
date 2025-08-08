# EZ Web Search MCP Server - Testing Guide

This document provides comprehensive testing instructions using official MCP tools and best practices.

## 🎯 Testing Overview

Our testing strategy follows MCP best practices with multiple testing levels:

1. **Official MCP Inspector** (Primary)
2. **Custom verification scripts** (Secondary)
3. **Manual protocol testing** (Development)
4. **Integration testing** (Client compatibility)

## 🔧 Official MCP Inspector Testing

The **MCP Inspector** is the official testing tool from the Model Context Protocol team.

### Prerequisites

- Node.js (for MCP Inspector)
- Built server binary (`./ez-web-search`)

### Quick Start

```bash
# Build the server
go build -o ez-web-search main.go

# Run guided testing
./test_with_inspector.sh
```

### CLI Mode Testing

Perfect for automation, CI/CD, and quick verification:

```bash
# 1. List all available tools
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/list

# 2. Test connectivity
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name ping

# 3. Test web search (basic)
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="Go programming tutorial"

# 4. Test web search with intent analysis
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="MCP testing best practices" --tool-arg search_intent=true

# 5. Test error handling (missing required parameter)
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search
```

### UI Mode Testing

Interactive web interface for development and debugging:

```bash
# Start Inspector UI (opens browser automatically)
npx @modelcontextprotocol/inspector ./ez-web-search
```

**Features:**
- Visual tool inspection
- Form-based parameter input
- Real-time response viewing
- Server log monitoring
- Configuration export for other clients

**Access:** http://localhost:6274

## 📋 Test Cases

### 1. Protocol Compliance Tests

#### Initialize Connection
```bash
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method initialize
```

**Expected:** Server info with protocol version 2024-11-05

#### Tool Discovery
```bash
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/list
```

**Expected:** Two tools (ping, web_search) with proper schemas

### 2. Functional Tests

#### Ping Tool
```bash
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name ping
```

**Expected:** `{"content": [{"type": "text", "text": "pong"}]}`

#### Web Search - Basic
```bash
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="Go programming"
```

**Expected:** Search results with titles, URLs, summaries

#### Web Search - With Intent
```bash
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="MCP testing" --tool-arg search_intent=true
```

**Expected:** Search results + intent analysis section

### 3. Error Handling Tests

#### Missing Required Parameter
```bash
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search
```

**Expected:** Error message about missing query parameter

#### Invalid Tool Name
```bash
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name invalid_tool
```

**Expected:** Tool not found error

## 🧪 Custom Testing Scripts

### Basic Verification
```bash
go run verify.go
```

**Tests:**
- Binary compilation
- Server startup
- Basic connectivity

### Interactive Python Testing
```bash
python3 interactive_test.py "your search query"
```

**Tests:**
- Full protocol flow
- Multiple tool calls
- Response parsing

## 📊 Test Results Validation

### Expected Tool Schema
```json
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
```

### Expected Search Response Format
```
Search Results for: [query]
Request ID: [request_id]

Search Intent Analysis: (if enabled)
- Query: [analyzed_query]
  Intent: [intent_type]
  Keywords: [extracted_keywords]

Search Results:
1. [title]
   URL: [url]
   Summary: [content_summary]
   Published: [date]
   Source: [source_ref]
...
```

## 🔗 Integration Testing

### Claude Desktop Integration
1. Export config using Inspector UI
2. Add to Claude Desktop's `mcp.json`
3. Test in Claude interface

### Cursor Integration
1. Use Inspector's "Servers File" export
2. Configure in Cursor settings
3. Verify tool availability

### VS Code Integration
1. Export configuration
2. Add to VS Code MCP settings
3. Test with GitHub Copilot

## 🚀 CI/CD Testing

For automated testing in CI/CD pipelines:

```bash
#!/bin/bash
# ci-test.sh

# Build
go build -o ez-web-search main.go

# Basic verification
go run verify.go

# Protocol compliance
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/list

# Functional tests
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name ping

# Integration test
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="test query"
```

## 🐛 Debugging

### Common Issues

1. **Server won't start**
   - Check Go version (1.21+)
   - Verify dependencies: `go mod tidy`

2. **Inspector connection fails**
   - Ensure server binary exists
   - Check file permissions: `chmod +x ez-web-search`

3. **Web search fails**
   - Verify API token
   - Check network connectivity
   - Review API rate limits

### Debug Mode

Enable verbose logging:
```bash
DEBUG=true ./ez-web-search
```

## 📈 Performance Testing

### Response Time Testing
```bash
time npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="performance test"
```

### Concurrent Testing
```bash
# Run multiple searches simultaneously
for i in {1..5}; do
  npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="test $i" &
done
wait
```

## ✅ Test Checklist

- [ ] Server builds successfully
- [ ] Basic verification passes
- [ ] Inspector CLI tests pass
- [ ] Inspector UI loads correctly
- [ ] All tools discovered properly
- [ ] Ping tool responds correctly
- [ ] Web search returns results
- [ ] Search intent analysis works
- [ ] Error handling works properly
- [ ] Configuration export works
- [ ] Integration with target client successful

## 📚 Additional Resources

- [MCP Inspector Documentation](https://modelcontextprotocol.io/docs/tools/inspector)
- [MCP Debugging Guide](https://modelcontextprotocol.io/docs/tools/debugging)
- [MCP Protocol Specification](https://modelcontextprotocol.io/specification)
- [BigModel Web Search API](https://docs.bigmodel.cn/api-reference/工具-api/网络搜索)

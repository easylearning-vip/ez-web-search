# EZ Web Search & Fetch MCP Server - Feature Overview

## 🚀 Complete AI Web Interaction Toolkit

This MCP server provides a comprehensive set of web interaction tools for AI applications, combining search and content extraction capabilities in a single, production-ready solution.

## 🛠️ Available Tools

### 1. 🔍 Web Search Tool
**Purpose**: Search the web using professional search APIs

**Capabilities**:
- Real-time web search using BigModel API
- Search intent analysis and keyword extraction
- Structured result formatting
- Multiple search result sources
- Publication date and source tracking

**Parameters**:
- `query` (required): Search query string
- `search_intent` (optional): Enable intent analysis

**Example Output**:
```
Search Results for: Go programming tutorial
Request ID: 202508080750181e139ad01ef440a4

Search Intent Analysis:
- Query: Go programming tutorial
  Intent: SEARCH_ALWAYS
  Keywords: go programming 教程

Search Results:
1. Ultimate Go Programming, Second Edition
   URL: https://zhuanlan.zhihu.com/p/617023696
   Summary: Ultimate Go Programming LiveLessons，第二版...
   Published: 2023-03-30
   Source: ref_1
...
```

### 2. 📄 Web Fetch Tool
**Purpose**: Extract content and metadata from any web page

**Capabilities**:
- Intelligent content extraction from HTML
- Metadata extraction (title, description, author, keywords, language)
- Link and image URL extraction
- Relative to absolute URL conversion
- Content cleaning and formatting
- Security validation (HTTP/HTTPS only)
- Configurable output options

**Parameters**:
- `url` (required): Target web page URL
- `include_links` (optional): Include extracted links
- `include_images` (optional): Include extracted images

**Example Output**:
```
Web Page Content for: https://example.com
Status Code: 200
Content Type: text/html

Title: Example Domain

Description: This domain is for use in illustrative examples

Author: IANA
Language: en
Keywords: example, domain, documentation

Content:
This domain is for use in illustrative examples in documents. 
You may use this domain in literature without prior coordination 
or asking for permission.

Links (5 found):
- https://www.iana.org/domains/example
- https://tools.ietf.org/html/rfc2606
...

Images (2 found):
- https://example.com/logo.png
- https://example.com/banner.jpg
```

### 3. 🏓 Ping Tool
**Purpose**: Test server connectivity and responsiveness

**Capabilities**:
- Simple connectivity verification
- Server health check
- MCP protocol validation

**Parameters**: None

**Example Output**:
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

## 🎯 Use Cases

### For AI Applications
- **Research Assistant**: Search for information and fetch detailed content
- **Content Analysis**: Extract and analyze web page content
- **Link Verification**: Check if URLs are accessible and extract metadata
- **Data Collection**: Gather structured information from web sources

### For Development
- **API Testing**: Verify web services and endpoints
- **Content Migration**: Extract content from existing websites
- **SEO Analysis**: Extract metadata and content structure
- **Link Checking**: Validate external references

## 🔧 Technical Features

### Content Extraction Intelligence
- **Multi-selector Strategy**: Tries multiple CSS selectors to find main content
- **Content Quality Filtering**: Only includes substantial text blocks (>50 chars)
- **Automatic Cleanup**: Removes excessive whitespace and formatting
- **Length Management**: Limits content to prevent excessive data (5000 chars max)

### URL Handling
- **Security Validation**: Only allows HTTP/HTTPS protocols
- **Relative URL Resolution**: Converts relative URLs to absolute
- **Error Handling**: Graceful handling of invalid or inaccessible URLs

### Performance Optimization
- **Timeout Management**: 30-second timeout for all requests
- **Resource Limiting**: Limits links (50 max) and images (20 max)
- **Memory Efficiency**: Streaming HTML parsing with goquery

### Protocol Compliance
- **MCP 2024-11-05**: Full compliance with latest protocol version
- **JSON-RPC 2.0**: Standard message format
- **Error Handling**: Proper error responses and status codes
- **Schema Validation**: Complete tool schema definitions

## 🧪 Testing & Validation

### Official MCP Inspector Support
- ✅ CLI mode testing for automation
- ✅ UI mode testing for development
- ✅ Full protocol compliance verification
- ✅ Tool discovery and schema validation

### Comprehensive Test Coverage
- ✅ All three tools functional
- ✅ Parameter validation
- ✅ Error handling scenarios
- ✅ Real API integration
- ✅ Content extraction accuracy
- ✅ URL security validation

## 🔒 Security Features

### Input Validation
- URL scheme validation (HTTP/HTTPS only)
- Parameter type checking
- Content length limiting
- Malformed URL handling

### Safe Defaults
- Conservative timeout settings
- Resource usage limits
- Error message sanitization
- No execution of JavaScript content

## 📊 Performance Characteristics

### Response Times
- **Ping**: < 10ms
- **Web Search**: 1-3 seconds (API dependent)
- **Web Fetch**: 1-5 seconds (page size dependent)

### Resource Usage
- **Memory**: Minimal footprint with streaming parsing
- **CPU**: Efficient HTML processing with goquery
- **Network**: Optimized HTTP client with connection reuse

## 🚀 Getting Started

```bash
# Build the server
go build -o ez-web-search main.go

# Test all tools
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/list

# Test web search
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="AI development"

# Test web fetch
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_fetch --tool-arg url="https://example.com"

# Start interactive UI
npx @modelcontextprotocol/inspector ./ez-web-search
```

## 🎉 Summary

This MCP server provides a complete web interaction toolkit for AI applications, combining:

- **Professional web search** with intent analysis
- **Intelligent content extraction** from any web page  
- **Robust error handling** and security validation
- **Full MCP protocol compliance** with official tool testing
- **Production-ready performance** with optimized resource usage

Perfect for AI assistants, research tools, content analysis applications, and any system requiring comprehensive web interaction capabilities.

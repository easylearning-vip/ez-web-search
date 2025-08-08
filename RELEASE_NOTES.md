# EZ Web Search & Fetch MCP Server - Release Notes

## v1.0.0 - First Stable Release 🎉

**Release Date**: August 8, 2025

### 🚀 What's New

This is the first stable release of EZ Web Search & Fetch MCP Server - a complete, enterprise-grade web interaction toolkit for AI applications.

### ✨ Key Features

#### 🔍 Web Search Tool
- Professional web search using BigModel API
- Search intent analysis and keyword extraction
- Structured result formatting with metadata
- Real-time search with comprehensive results

#### 📄 Web Fetch Tool
- Intelligent content extraction from any web page
- Metadata extraction (title, description, author, keywords, language)
- Link and image URL extraction with absolute URL conversion
- Content cleaning and formatting
- Configurable output options

#### 🛡️ Advanced Anti-Bot Protection
- User agent rotation (12+ realistic browser UAs)
- Request header spoofing (complete browser headers)
- Random delays between requests (1-3s configurable)
- Intelligent retry logic with exponential backoff
- Rate limiting and blocking detection
- WAF/CDN bypass mechanisms

#### 🏗️ Enterprise Architecture
- Modular design with clean separation of concerns
- Environment-based configuration management
- Comprehensive error handling and logging
- Performance optimizations and resource limits
- Security best practices implementation

#### 🤖 Claude Code CLI Integration
- Ready-to-use configuration templates
- Automated setup scripts
- Comprehensive usage examples
- Troubleshooting guides

### 📦 Installation

#### One-Click Installation (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/easylearning-vip/ez-web-search/main/install.sh | bash
```

#### Manual Installation
1. Download the appropriate binary for your platform from the [releases page](https://github.com/easylearning-vip/ez-web-search/releases/latest)
2. Make it executable: `chmod +x ez-web-search-v2`
3. Run the setup: `./setup-claude-cli.sh`

### 🖥️ Supported Platforms

- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

### 📋 Usage Examples

#### Web Search
```bash
# In Claude Code CLI
> Search for "Go web scraping best practices"
> Find recent articles about MCP protocol implementation
```

#### Web Fetch
```bash
# In Claude Code CLI
> Fetch content from https://pkg.go.dev/github.com/PuerkitoBio/goquery
> Get documentation from https://modelcontextprotocol.io with links included
```

#### Combined Workflows
```bash
# In Claude Code CLI
> Search for "Go HTTP client tutorials", then fetch content from the top 3 results
> Find MCP server examples, fetch their GitHub repos, and analyze the code
```

### 🔧 Configuration

The server supports extensive configuration via environment variables:

```bash
# API Configuration
BIGMODEL_TOKEN="your_api_token"
BIGMODEL_TIMEOUT="30s"

# Anti-Bot Settings
WEBFETCH_USER_AGENT_ROTATE=true
WEBFETCH_DELAY_MIN="1s"
WEBFETCH_DELAY_MAX="3s"

# Resource Limits
WEBFETCH_MAX_CONTENT_SIZE=5000
WEBFETCH_MAX_LINKS=50
WEBFETCH_MAX_IMAGES=20
```

### 🧪 Testing

The release includes comprehensive testing tools:

```bash
# Test all tools with MCP Inspector
make test-all-tools

# Test individual components
make test-ping
make test-search
make test-fetch

# Start interactive UI
make inspector-ui
```

### 📚 Documentation

- **README**: Complete setup and usage guide
- **FEATURES**: Detailed feature overview
- **ENTERPRISE_FEATURES**: Enterprise architecture details
- **TESTING**: Comprehensive testing guide

### 🔒 Security

- Environment-based configuration (no hardcoded secrets)
- Input validation and sanitization
- Resource limits and timeout management
- Secure error handling

### 🚀 Performance

- Efficient HTTP client with connection reuse
- Streaming HTML parsing for memory efficiency
- Configurable resource limits
- Optimized binary size (6-7MB)

### 🐛 Known Issues

None reported for this release.

### 📈 What's Next

- Additional search engines support
- Enhanced content extraction algorithms
- Caching mechanisms
- Monitoring and metrics
- Plugin system for extensibility

### 🙏 Acknowledgments

- Built with [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - Official MCP Go library
- HTML parsing powered by [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)
- Web search provided by BigModel API

### 📞 Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/easylearning-vip/ez-web-search/issues)
- **Documentation**: [GitHub Repository](https://github.com/easylearning-vip/ez-web-search)
- **Discussions**: [GitHub Discussions](https://github.com/easylearning-vip/ez-web-search/discussions)

---

**Full Changelog**: https://github.com/easylearning-vip/ez-web-search/commits/v1.0.0

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.1.1] - 2025-08-25

### 🧪 Testing & Quality Improvements

#### Enhanced Testing Suite
- **Complete Test Script**: Added comprehensive `test_all_features.sh` script
  - Tests all 3 tools (ping, ez_web_search, ez_web_fetch)
  - Validates all 4 search engines (search_std, search_pro, search_pro_sogou, search_pro_quark)
  - Tests search intent analysis functionality
  - Tests web fetch with links and images options
  - Validates error handling scenarios
  - Provides detailed test results and status reporting

#### Documentation Updates
- **Enhanced Testing Documentation**: Updated TESTING.md with latest test procedures
- **Feature Documentation**: Improved FEATURES.md with current capabilities
- **Claude Integration Guide**: Added CLAUDE.md for AI assistant integration
- **Project Summary**: Updated PROJECT_SUMMARY.md with latest status

#### Quality Assurance
- ✅ **MCP Inspector Compliance**: All tests pass with official MCP Inspector
- ✅ **Multi-Engine Validation**: All 4 search engines tested and working
- ✅ **Error Handling**: Proper error responses validated
- ✅ **Production Ready**: Complete test coverage for production deployment

### 🔧 Configuration & Environment
- **Environment File**: Added `.env` template for easy configuration
- **Token Management**: Improved API token handling and validation
- **Build Verification**: Enhanced build process validation

---

## [v1.1.0] - 2025-08-08

### 🎉 Major Features Added

#### 🔍 Configurable Search Engines
- **Multiple Search Engine Support**: Added support for 4 different search engines
  - `search_std`: 智谱基础版搜索引擎 (default)
  - `search_pro`: 智谱高阶版搜索引擎
  - `search_pro_sogou`: 搜狗搜索
  - `search_pro_quark`: 夸克搜索

- **Flexible Configuration Options**:
  - Environment variable: `BIGMODEL_SEARCH_ENGINE` (sets default)
  - Tool parameter: `search_engine` (runtime override)
  - Parameter validation for supported engines
  - Backward compatibility maintained

#### 📚 Multi-Language Documentation
- **Chinese README**: Comprehensive Chinese documentation (`README-zh.md`)
- **Language Selection**: Added language links in both README files
- **Enhanced Documentation**: Updated configuration examples and usage guides

### 🏗️ Architecture Improvements

#### Enhanced Configuration Management
- **BigModelConfig**: Added `SearchEngine` field
- **WebSearchOptions**: Added `SearchEngine` field for runtime configuration
- **Environment Variables**: Enhanced `.env.example` with search engine options

#### Service Layer Updates
- **WebSearchService**: Enhanced to use configurable search engines
- **FormatSearchResponse**: Updated to display selected search engine
- **MCP Handlers**: Added search engine parameter support and validation

#### Code Structure Cleanup
- **Modular Architecture**: Removed conflicting `main.go`, using `cmd/server/main.go`
- **Type Safety**: Enhanced type definitions with search engine support
- **Validation**: Added search engine parameter validation

### 🧪 Testing & Quality

#### MCP Protocol Compliance
- ✅ Updated tool schema with `search_engine` parameter
- ✅ MCP Inspector validation passed
- ✅ Tool discovery and execution verified

#### Multi-Platform Support
- ✅ Linux (amd64, arm64)
- ✅ macOS (amd64, arm64) 
- ✅ Windows (amd64)
- ✅ Binary size optimized (6-7MB)

### 📦 Distribution & Installation

#### Release Artifacts
- **Multi-Platform Binaries**: Built for 5 platforms
- **Checksums**: SHA256 verification files
- **Release Notes**: Comprehensive release documentation

#### Installation Methods
- **One-Click Install**: `curl -fsSL https://raw.githubusercontent.com/easylearning-vip/ez-web-search/main/install.sh | bash`
- **Manual Install**: Platform-specific binary downloads
- **Source Build**: Enhanced build system with release automation

### 🔧 Configuration Examples

#### Environment Configuration
```bash
# Default search engine (optional)
BIGMODEL_SEARCH_ENGINE="search_std"

# API Configuration
BIGMODEL_TOKEN="your_api_token"
BIGMODEL_TIMEOUT="30s"

# Anti-Bot Settings
WEBFETCH_USER_AGENT_ROTATE=true
WEBFETCH_DELAY_MIN="1s"
WEBFETCH_DELAY_MAX="3s"
```

#### Tool Usage Examples
```bash
# Use default search engine
> Search for "Go programming tutorials"

# Override with specific search engine
> Search for "web scraping" using search_pro engine

# Use Sogou search
> Use search_pro_sogou to find "MCP protocol examples"

# Use Quark search with intent analysis
> Search with search_pro_quark for "AI development" with search intent analysis
```

### 🌐 Multi-Language Support

#### Documentation Languages
- **English**: `README.md` - Complete English documentation
- **中文**: `README-zh.md` - Comprehensive Chinese documentation

#### Localized Content
- **Feature Descriptions**: Translated feature explanations
- **Configuration Guides**: Localized setup instructions
- **Usage Examples**: Language-appropriate examples
- **Installation Instructions**: Multi-language installation guides

### ⚠️ Breaking Changes

#### API Changes
- **FormatSearchResponse**: Method signature updated to include `searchEngine` parameter
- **WebSearchOptions**: Struct enhanced with `SearchEngine` field

#### Migration Guide
- **Existing Configurations**: No changes required, backward compatible
- **Custom Implementations**: Update method calls if using `FormatSearchResponse` directly
- **Environment Variables**: Optional `BIGMODEL_SEARCH_ENGINE` can be added

### 🔄 Backward Compatibility

#### Maintained Compatibility
- **Default Behavior**: Uses `search_std` if no engine specified
- **Existing Configurations**: All existing `.env` files continue to work
- **Tool Interface**: Previous tool calls without `search_engine` parameter work unchanged
- **API Responses**: Response format remains consistent

### 📈 Performance & Reliability

#### Optimizations
- **Search Engine Selection**: Efficient validation and fallback logic
- **Memory Usage**: Optimized struct layouts and string handling
- **Binary Size**: Maintained compact binary size despite new features
- **Startup Time**: Fast initialization with enhanced configuration loading

#### Reliability Improvements
- **Input Validation**: Enhanced parameter validation for search engines
- **Error Handling**: Improved error messages for invalid search engines
- **Fallback Logic**: Graceful degradation to default search engine
- **Configuration Validation**: Startup-time configuration verification

### 🚀 Future Roadmap

#### Planned Enhancements
- Additional search engine integrations
- Search result caching mechanisms
- Advanced search filtering options
- Performance metrics and monitoring
- Plugin system for custom search engines

---

## [v1.0.0] - 2025-08-08

### 🎉 Initial Release

#### Core Features
- **Web Search Tool**: BigModel API integration with anti-bot protection
- **Web Fetch Tool**: Intelligent content extraction from web pages
- **Enterprise Architecture**: Modular, scalable, production-ready design
- **Claude Code CLI Integration**: Seamless AI assistant integration
- **Multi-Platform Support**: Linux, macOS, Windows binaries

#### Anti-Bot Protection
- User agent rotation (12+ realistic browser UAs)
- Request header spoofing and randomization
- Intelligent retry logic with exponential backoff
- Rate limiting and blocking detection
- WAF/CDN bypass mechanisms

#### Documentation & Distribution
- Comprehensive README and documentation
- One-click installation script
- Multi-platform binary releases
- MCP Inspector compliance verification
- Enterprise deployment guides

---

**Repository**: https://github.com/easylearning-vip/ez-web-search  
**Documentation**: [English](README.md) | [中文](README-zh.md)  
**Latest Release**: [v1.1.0](https://github.com/easylearning-vip/ez-web-search/releases/latest)

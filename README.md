# EZ Web Search & Fetch MCP Server

**English** | [中文](README-zh.md)

A complete, enterprise-grade Web Search and Fetch MCP (Model Context Protocol) server implemented in Go, providing comprehensive web interaction capabilities for AI applications with advanced anti-bot protection and configurable search engines.

## 🚀 Features

- **Web Search Tool**: Search the web using BigModel's Web Search API with configurable search engines
  - **search_std**: 智谱基础版搜索引擎 (default)
  - **search_pro**: 智谱高阶版搜索引擎
  - **search_pro_sogou**: 搜狗搜索
  - **search_pro_quark**: 夸克搜索
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

### One-Click Installation (Recommended)

```bash
# Install and configure automatically
curl -fsSL https://raw.githubusercontent.com/easylearning-vip/ez-web-search/main/install.sh | bash
```

This script will:
- Download the latest release for your platform
- Install the binary to `~/.local/bin`
- Configure Claude Code CLI automatically
- Set up your BigModel API token
- Test the installation


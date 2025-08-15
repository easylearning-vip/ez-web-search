# EZ Web Search & Fetch MCP 服务器

[English](README.md) | **中文**

一个企业级的网络搜索和内容获取 MCP (Model Context Protocol) 服务器，专为 AI 应用设计，具备强大的防反爬虫机制。

## 🌟 主要特性

### 🔍 网络搜索工具
- **多搜索引擎支持**: 智谱基础版、高阶版、搜狗、夸克搜索
- **智能搜索意图分析**: 关键词提取和搜索意图识别
- **结构化结果**: 包含元数据的格式化搜索结果
- **实时搜索**: 快速响应的网络搜索功能

### 📄 网页内容获取工具
- **智能内容提取**: 从任何网页提取主要内容
- **元数据提取**: 标题、描述、作者、关键词、语言等
- **链接和图片提取**: 自动转换为绝对URL
- **内容清理和格式化**: 智能去除噪音内容
- **可配置输出选项**: 灵活的内容包含设置

### 🛡️ 高级防反爬虫保护
- **用户代理轮换**: 12+个真实浏览器用户代理池
- **请求头伪装**: 完整的浏览器请求头模拟
- **随机延迟**: 1-3秒可配置的请求间隔
- **智能重试逻辑**: 指数退避重试机制
- **速率限制检测**: 自动检测429、503状态码
- **WAF/CDN绕过**: 识别和绕过常见阻断模式


## 🚀 快速开始

### 一键安装（推荐）

```bash
# 自动安装和配置
curl -fsSL https://raw.githubusercontent.com/easylearning-vip/ez-web-search/main/install.sh | bash
```

此脚本将：
- 下载适合您平台的最新版本
- 安装二进制文件到 `~/.local/bin`
- 自动配置 Claude Code CLI
- 设置您的 BigModel API token
- 测试安装

## 配置

### Claude CLI

```
claude mcp add web-search -- ez-web-search --token your-token
```


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
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_fetch --tool-arg url="https://www.easylearning.vip"

# Test web fetch with links and images
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_fetch --tool-arg url="https://www.easylearning.vip" --tool-arg include_links=true --tool-arg include_images=true

# Test with search intent analysis
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --tool-arg query="MCP testing" --tool-arg search_intent=true
```

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

### 🏗️ 企业级架构
- **模块化设计**: 清晰的关注点分离
- **环境配置管理**: 无硬编码敏感信息
- **全面错误处理**: 健壮的错误处理和日志记录
- **性能优化**: 连接复用和资源限制
- **安全最佳实践**: 遵循安全开发规范

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

### 手动安装

#### 从发布版本安装（推荐）

1. **下载最新版本**:
   ```bash
   # 访问发布页面下载适合您平台的版本
   # https://github.com/easylearning-vip/ez-web-search/releases/latest
   
   # 或使用 curl（替换为您的平台）
   curl -L -o ez-web-search \
     "https://github.com/easylearning-vip/ez-web-search/releases/latest/download/ez-web-search_linux_amd64"
   
   chmod +x ez-web-search
   ```

2. **配置 Claude Code CLI**:
   ```bash
   ./setup-claude-cli.sh
   ```

#### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/easylearning-vip/ez-web-search.git
cd ez-web-search

# 构建服务器
make build

# 使用默认配置运行
make run

# 运行所有测试
make test-all-tools

# 启动 Inspector UI
make inspector-ui
```

## 🔧 配置

### 环境变量配置

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件设置您的配置
# BIGMODEL_TOKEN="your_actual_bigmodel_api_token"
# BIGMODEL_SEARCH_ENGINE="search_std"
# WEBFETCH_USER_AGENT_ROTATE=true
# WEBFETCH_DELAY_MIN="1s"
# WEBFETCH_DELAY_MAX="3s"

# 使用环境变量运行
make dev
```

### 搜索引擎选项

- **search_std**: 智谱基础版搜索引擎（默认）
- **search_pro**: 智谱高阶版搜索引擎
- **search_pro_sogou**: 搜狗搜索
- **search_pro_quark**: 夸克搜索

可以通过以下方式设置搜索引擎：
1. 环境变量：`BIGMODEL_SEARCH_ENGINE=search_pro`
2. 工具调用参数：在调用时指定 `search_engine` 参数

## 🤖 Claude Code CLI 集成

### 快速设置

**自动化设置（推荐）**:
```bash
# 运行自动化设置脚本
./setup-claude-cli.sh
```

**手动设置**:

1. **构建服务器**:
   ```bash
   make build
   ```

2. **复制配置模板**:
   ```bash
   # 创建 Claude Code CLI MCP 配置目录
   mkdir -p ~/.claude
   
   # 复制并自定义配置模板
   cp claude-mcp-config.json ~/.claude/mcp_settings.json
   ```

3. **更新配置**:
   ```bash
   # 更新二进制文件路径
   PWD_PATH=$(pwd)
   sed -i "s|/path/to/ez-web-search|$PWD_PATH/ez-web-search|g" ~/.claude/mcp_settings.json
   
   # 设置您的 BigModel API token
   sed -i 's/your_bigmodel_api_token_here/YOUR_ACTUAL_TOKEN/g' ~/.claude/mcp_settings.json
   ```

4. **启动 Claude Code CLI**:
   ```bash
   claude
   ```

### 使用示例

#### 网络搜索
```bash
# 在 Claude Code CLI 中
> 搜索 "Go 网络爬虫最佳实践"
> 使用搜狗搜索引擎查找 "MCP 协议实现"
> 搜索 "Go HTTP 客户端教程" 并启用搜索意图分析
```

#### 网页获取
```bash
# 在 Claude Code CLI 中
> 获取 https://pkg.go.dev/github.com/PuerkitoBio/goquery 的内容并总结主要特性
> 从 https://modelcontextprotocol.io 获取文档并提取关键概念
> 获取 https://example.com 的内容，包含链接和图片
```

#### 组合工作流
```bash
# 在 Claude Code CLI 中
> 搜索 "Go HTTP 客户端最佳实践"，然后获取前3个结果的内容并创建综合指南
> 查找最新的 Go 网络爬虫库，获取它们的文档，并比较功能特性
> 搜索 MCP 服务器示例，获取 GitHub 仓库，并分析代码结构
```

## 🧪 测试

### MCP Inspector 测试

```bash
# 使用 MCP Inspector 测试所有工具
make test-all-tools

# 测试单个组件
make test-ping
make test-search
make test-fetch

# 启动交互式 UI
make inspector-ui
```

### 工具测试示例

```bash
# 测试 ping 工具
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name ping

# 测试网络搜索
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_search --args '{"query": "Go 编程教程", "search_engine": "search_pro"}'

# 测试网页获取
npx @modelcontextprotocol/inspector --cli ./ez-web-search --method tools/call --tool-name web_fetch --args '{"url": "https://example.com", "include_links": true}'
```

## 📊 支持的平台

- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

## 🔒 安全特性

- 基于环境变量的配置（无硬编码密钥）
- 输入验证和清理
- 资源限制和超时管理
- 安全的错误处理

## 🚀 性能特性

- 高效的 HTTP 客户端连接复用
- 流式 HTML 解析，内存高效
- 可配置的资源限制
- 优化的二进制大小（6-7MB）

## 📚 文档

- **README**: 完整的设置和使用指南
- **FEATURES**: 详细的功能概述
- **ENTERPRISE_FEATURES**: 企业架构详情
- **TESTING**: 综合测试指南

## 🤝 贡献

欢迎贡献！请查看我们的贡献指南。

## 📄 许可证

本项目采用 MIT 许可证。

## 📞 支持

- **GitHub Issues**: [报告错误或请求功能](https://github.com/easylearning-vip/ez-web-search/issues)
- **文档**: [GitHub 仓库](https://github.com/easylearning-vip/ez-web-search)
- **讨论**: [GitHub 讨论](https://github.com/easylearning-vip/ez-web-search/discussions)

---

**项目状态**: ✅ **完成并可用于生产**  
**最后更新**: 2025年8月8日

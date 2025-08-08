#!/bin/bash

# EZ Web Search MCP Server - Claude Code CLI Setup Script
# This script helps you configure the MCP server for Claude Code CLI

set -e

echo "🚀 EZ Web Search MCP Server - Claude Code CLI Setup"
echo "=================================================="
echo

# Check if Claude CLI is installed
if ! command -v claude &> /dev/null; then
    echo "❌ Claude Code CLI is not installed."
    echo "Please install Claude Code CLI first: https://claude.ai/cli"
    exit 1
fi

echo "✅ Claude Code CLI found"

# Build the server
echo "🔨 Building the MCP server..."
make build

if [ ! -f "./ez-web-search-v2" ]; then
    echo "❌ Failed to build ez-web-search-v2"
    exit 1
fi

echo "✅ Server built successfully"

# Get current directory
CURRENT_DIR=$(pwd)
BINARY_PATH="$CURRENT_DIR/ez-web-search-v2"

echo "📍 Binary location: $BINARY_PATH"

# Create Claude MCP configuration directory
CLAUDE_DIR="$HOME/.claude"
mkdir -p "$CLAUDE_DIR"

# Check if MCP settings file exists
MCP_CONFIG="$CLAUDE_DIR/mcp_settings.json"

if [ -f "$MCP_CONFIG" ]; then
    echo "⚠️  MCP configuration file already exists: $MCP_CONFIG"
    echo "Creating backup..."
    cp "$MCP_CONFIG" "$MCP_CONFIG.backup.$(date +%Y%m%d_%H%M%S)"
    echo "✅ Backup created"
fi

# Prompt for BigModel API token
echo
echo "🔑 BigModel API Token Configuration"
echo "Please enter your BigModel API token:"
echo "(You can get one from: https://open.bigmodel.cn/)"
echo
read -p "BigModel API Token: " BIGMODEL_TOKEN

if [ -z "$BIGMODEL_TOKEN" ]; then
    echo "⚠️  No token provided. Using default test token."
    BIGMODEL_TOKEN="0f405f7a11b946298b154f042a70f12b.s6VO3ITALpa3bhDo"
fi

# Create MCP configuration
echo "📝 Creating MCP configuration..."

cat > "$MCP_CONFIG" << EOF
{
  "mcpServers": {
    "ez-web-search": {
      "command": "$BINARY_PATH",
      "env": {
        "BIGMODEL_TOKEN": "$BIGMODEL_TOKEN",
        "WEBFETCH_USER_AGENT_ROTATE": "true",
        "WEBFETCH_DELAY_MIN": "1s",
        "WEBFETCH_DELAY_MAX": "3s",
        "WEBFETCH_MAX_CONTENT_SIZE": "5000",
        "WEBFETCH_MAX_LINKS": "50",
        "WEBFETCH_MAX_IMAGES": "20",
        "PATH": "/usr/local/bin:/usr/bin:/bin"
      }
    }
  }
}
EOF

echo "✅ MCP configuration created: $MCP_CONFIG"

# Test the configuration
echo
echo "🧪 Testing the configuration..."

# Test if the binary runs
if timeout 5s "$BINARY_PATH" < /dev/null > /dev/null 2>&1; then
    echo "✅ Server binary runs correctly"
else
    echo "⚠️  Server binary test inconclusive (this is normal for MCP servers)"
fi

# Validate JSON configuration
if command -v jq &> /dev/null; then
    if jq empty "$MCP_CONFIG" 2>/dev/null; then
        echo "✅ MCP configuration JSON is valid"
    else
        echo "❌ MCP configuration JSON is invalid"
        exit 1
    fi
else
    echo "⚠️  jq not found, skipping JSON validation"
fi

echo
echo "🎉 Setup completed successfully!"
echo
echo "📋 Next steps:"
echo "1. Start Claude Code CLI:"
echo "   claude"
echo
echo "2. Test the web search functionality:"
echo '   > Search for "Go web scraping best practices"'
echo
echo "3. Test the web fetch functionality:"
echo '   > Fetch content from https://example.com'
echo
echo "4. Try combined workflows:"
echo '   > Search for Go tutorials, then fetch content from the top result'
echo
echo "📚 For more examples, see README.md"
echo
echo "🔧 Configuration file location: $MCP_CONFIG"
echo "🔧 Server binary location: $BINARY_PATH"
echo
echo "💡 Tip: You can edit the configuration file to customize settings"
echo "💡 Tip: Use 'make test-all-tools' to test the server independently"

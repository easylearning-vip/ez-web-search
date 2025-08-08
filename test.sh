#!/bin/bash

# EZ Web Search MCP Server Test Script

echo "EZ Web Search MCP Server Test"
echo "============================="
echo

# Check if the binary exists
if [ ! -f "./ez-web-search" ]; then
    echo "Building the server..."
    go build -o ez-web-search main.go
    if [ $? -ne 0 ]; then
        echo "Failed to build the server"
        exit 1
    fi
fi

echo "Testing the MCP server..."
echo

# Test with a sample query
QUERY="Go programming language tutorial"
if [ $# -gt 0 ]; then
    QUERY="$1"
fi

echo "Running test with query: '$QUERY'"
echo

# Run the test client
go run test_client.go "$QUERY"

echo
echo "Test completed!"

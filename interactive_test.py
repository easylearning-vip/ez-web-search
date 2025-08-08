#!/usr/bin/env python3
"""
Interactive test client for EZ Web Search MCP Server
"""

import json
import subprocess
import sys
import time

def send_message(process, message):
    """Send a JSON-RPC message to the MCP server"""
    json_msg = json.dumps(message)
    print(f"→ Sending: {json_msg}")
    process.stdin.write(json_msg + '\n')
    process.stdin.flush()
    
    # Give some time for response
    time.sleep(1)

def main():
    if len(sys.argv) < 2:
        query = "Go programming tutorial"
        print(f"Using default query: {query}")
        print("Usage: python3 interactive_test.py '<your_search_query>'")
    else:
        query = sys.argv[1]
    
    print("EZ Web Search MCP Server Interactive Test")
    print("=========================================")
    print()
    
    # Start the MCP server
    print("Starting MCP server...")
    try:
        process = subprocess.Popen(
            ['./ez-web-search'],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        
        # Test 1: Initialize
        print("\n1. Initializing...")
        init_msg = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {
                    "name": "interactive-test",
                    "version": "1.0.0"
                }
            }
        }
        send_message(process, init_msg)
        
        # Test 2: List tools
        print("\n2. Listing tools...")
        list_msg = {
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list",
            "params": {}
        }
        send_message(process, list_msg)
        
        # Test 3: Ping
        print("\n3. Testing ping...")
        ping_msg = {
            "jsonrpc": "2.0",
            "id": 3,
            "method": "tools/call",
            "params": {
                "name": "ping",
                "arguments": {}
            }
        }
        send_message(process, ping_msg)
        
        # Test 4: Web search
        print(f"\n4. Testing web search with query: '{query}'...")
        search_msg = {
            "jsonrpc": "2.0",
            "id": 4,
            "method": "tools/call",
            "params": {
                "name": "web_search",
                "arguments": {
                    "query": query,
                    "search_intent": False
                }
            }
        }
        send_message(process, search_msg)
        
        # Test 5: Web search with intent
        print(f"\n5. Testing web search with intent analysis...")
        search_intent_msg = {
            "jsonrpc": "2.0",
            "id": 5,
            "method": "tools/call",
            "params": {
                "name": "web_search",
                "arguments": {
                    "query": query,
                    "search_intent": True
                }
            }
        }
        send_message(process, search_intent_msg)
        
        # Wait for all responses
        print("\nWaiting for responses...")
        time.sleep(5)
        
        # Read any available output
        try:
            stdout, stderr = process.communicate(timeout=2)
            if stdout:
                print("\n📤 Server responses:")
                print(stdout)
            if stderr:
                print("\n⚠️  Server errors:")
                print(stderr)
        except subprocess.TimeoutExpired:
            print("⏰ Timeout waiting for responses")
            process.kill()
            stdout, stderr = process.communicate()
            if stdout:
                print("\n📤 Partial server responses:")
                print(stdout)
        
        print("\n✅ Interactive test completed!")
        
    except FileNotFoundError:
        print("❌ Error: ./ez-web-search binary not found. Please run 'go build -o ez-web-search main.go' first.")
    except Exception as e:
        print(f"❌ Error: {e}")
    finally:
        if 'process' in locals():
            process.terminate()

if __name__ == "__main__":
    main()

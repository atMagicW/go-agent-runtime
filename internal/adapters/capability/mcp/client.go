package mcp

import (
	"context"
	"fmt"

	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// Client 是第一版 MCP 客户端 mock 实现
type Client struct {
}

// NewClient 创建 MCP 客户端
func NewClient() *Client {
	return &Client{}
}

// CallTool 调用远程 MCP Tool
func (c *Client) CallTool(_ context.Context, req ports.MCPCallRequest) (ports.MCPCallResponse, error) {
	// 第一版先返回 mock 数据。
	// 后续替换成真实 MCP SDK 调用时，只改这里。
	return ports.MCPCallResponse{
		Output: map[string]any{
			"server_name": req.ServerName,
			"tool_name":   req.ToolName,
			"result":      fmt.Sprintf("mcp tool %s executed on server %s", req.ToolName, req.ServerName),
		},
	}, nil
}

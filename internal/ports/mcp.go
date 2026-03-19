package ports

import "context"

// MCPCallRequest 表示一次 MCP Tool 调用请求
type MCPCallRequest struct {
	ServerName string
	ToolName   string
	Input      map[string]any
}

// MCPCallResponse 表示一次 MCP Tool 调用结果
type MCPCallResponse struct {
	Output map[string]any
}

// MCPClient 定义 MCP 客户端接口
type MCPClient interface {
	CallTool(ctx context.Context, req MCPCallRequest) (MCPCallResponse, error)
}

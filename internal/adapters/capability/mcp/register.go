package mcp

import "github.com/atMagicW/go-agent-runtime/internal/domain/capability"

// DefaultToolSpecs 返回默认内置的 MCP Tool 列表
func DefaultToolSpecs() []capability.MCPToolSpec {
	return []capability.MCPToolSpec{
		{
			Name:        "mcp_web_search",
			ServerName:  "search-server",
			RemoteTool:  "web_search",
			Description: "通过远程 MCP Server 执行 Web 搜索",
			Version:     "v1",
			Enabled:     true,
		},
		{
			Name:        "mcp_doc_lookup",
			ServerName:  "docs-server",
			RemoteTool:  "doc_lookup",
			Description: "通过远程 MCP Server 查询文档",
			Version:     "v1",
			Enabled:     true,
		},
	}
}

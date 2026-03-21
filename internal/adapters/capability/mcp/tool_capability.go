package mcp

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// ToolCapability 是把远程 MCP Tool 适配成统一 Capability 的实现
type ToolCapability struct {
	client ports.MCPClient
	spec   capability.MCPToolSpec
}

// NewToolCapability 创建一个 MCP Tool Capability
func NewToolCapability(client ports.MCPClient, spec capability.MCPToolSpec) *ToolCapability {
	return &ToolCapability{
		client: client,
		spec:   spec,
	}
}

// Descriptor 返回能力元信息
func (c *ToolCapability) Descriptor() capability.Descriptor {
	return capability.Descriptor{
		Name:        c.spec.Name,
		Kind:        capability.KindMCPTool,
		Description: c.spec.Description,
		Tags:        []string{"mcp", "remote_tool", c.spec.ServerName},
		Version:     c.spec.Version,
		Enabled:     c.spec.Enabled,
	}
}

// Invoke 执行 MCP Tool
func (c *ToolCapability) Invoke(ctx context.Context, input map[string]any) (capability.Result, error) {
	resp, err := c.client.CallTool(ctx, ports.MCPCallRequest{
		ServerName: c.spec.ServerName,
		ToolName:   c.spec.RemoteTool,
		Input:      input,
	})
	if err != nil {
		return capability.Result{
			Name:    c.spec.Name,
			Kind:    capability.KindMCPTool,
			Success: false,
			Error:   err.Error(),
		}, err
	}

	output := map[string]any{
		"capability_name":    c.spec.Name,
		"kind":               "mcp_tool",
		"server_name":        c.spec.ServerName,
		"server_description": c.spec.ServerDescription,
		"remote_tool":        c.spec.RemoteTool,
	}

	for k, v := range resp.Output {
		output[k] = v
	}

	return capability.Result{
		Name:    c.spec.Name,
		Kind:    capability.KindMCPTool,
		Success: true,
		Output:  output,
	}, nil
}

package app

import (
	"strings"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
)

// MCPService 提供 MCP 配置查询能力
type MCPService struct {
	cfg *cfg.CapabilitiesConfig
}

// NewMCPService 创建 MCPService
func NewMCPService(c *cfg.CapabilitiesConfig) *MCPService {
	return &MCPService{
		cfg: c,
	}
}

// ListServers 返回面向 API 的 MCP server 视图
func (s *MCPService) ListServers() []httpapi.MCPServerView {
	if s == nil || s.cfg == nil {
		return nil
	}

	out := make([]httpapi.MCPServerView, 0, len(s.cfg.MCPServers))
	for _, server := range s.cfg.MCPServers {
		mode := strings.TrimSpace(server.Mode)
		if mode == "" {
			mode = "mock"
		}

		tools := make([]httpapi.MCPToolView, 0, len(server.Tools))
		for _, tool := range server.Tools {
			tools = append(tools, httpapi.MCPToolView{
				Name:        tool.Name,
				RemoteTool:  tool.RemoteTool,
				Description: tool.Description,
				Enabled:     tool.Enabled,
			})
		}

		out = append(out, httpapi.MCPServerView{
			Name:        server.Name,
			Description: server.Description,
			Mode:        mode,
			BaseURL:     server.BaseURL,
			ToolPath:    server.ToolPath,
			TimeoutMS:   server.TimeoutMS,
			Enabled:     server.Enabled,
			Tools:       tools,
		})
	}

	return out
}

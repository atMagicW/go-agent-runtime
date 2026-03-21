package app

import cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"

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

// ListServers 返回所有 MCP server 配置
func (s *MCPService) ListServers() []cfg.MCPServerConfig {
	if s == nil || s.cfg == nil {
		return nil
	}

	out := make([]cfg.MCPServerConfig, len(s.cfg.MCPServers))
	copy(out, s.cfg.MCPServers)
	return out
}

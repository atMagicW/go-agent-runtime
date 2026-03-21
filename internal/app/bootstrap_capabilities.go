package app

import (
	capregistry "github.com/atMagicW/go-agent-runtime/internal/adapters/capability"
	mcpcap "github.com/atMagicW/go-agent-runtime/internal/adapters/capability/mcp"
	"github.com/atMagicW/go-agent-runtime/internal/adapters/capability/skills"
	"github.com/atMagicW/go-agent-runtime/internal/adapters/capability/tools"
	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// BuildCapabilityRegistry 从配置构建能力注册表
func BuildCapabilityRegistry(capCfg *cfg.CapabilitiesConfig, mcpClient ports.MCPClient) *capregistry.Registry {
	registry := capregistry.NewRegistry()

	if capCfg == nil {
		return registry
	}

	for _, item := range capCfg.Skills {
		if !item.Enabled {
			continue
		}
		switch item.Name {
		case "resume_analyzer":
			registry.MustRegister(skills.NewResumeAnalyzerSkill())
		}
	}

	for _, item := range capCfg.Tools {
		if !item.Enabled {
			continue
		}
		switch item.Name {
		case "keyword_extract_tool":
			registry.MustRegister(tools.NewKeywordExtractTool())
		}
	}

	for _, server := range capCfg.MCPServers {
		if !server.Enabled {
			continue
		}

		for _, tool := range server.Tools {
			if !tool.Enabled {
				continue
			}

			registry.MustRegister(mcpcap.NewToolCapability(mcpClient, capability.MCPToolSpec{
				Name:              tool.Name,
				ServerName:        server.Name,
				ServerDescription: server.Description,
				RemoteTool:        tool.RemoteTool,
				Description:       tool.Description,
				Version:           "v1",
				Enabled:           tool.Enabled,
			}))
		}
	}

	return registry
}

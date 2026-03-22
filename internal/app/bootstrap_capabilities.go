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
func BuildCapabilityRegistry(
	capCfg *cfg.CapabilitiesConfig,
	skillRegistry *SkillRegistry,
	mcpClient ports.MCPClient,
) *capregistry.Registry {
	registry := capregistry.NewRegistry()

	if capCfg == nil {
		return registry
	}

	for _, item := range capCfg.Skills {
		if !item.Enabled {
			continue
		}

		switch item.Name {
		case capability.CapabilityResumeAnalyzer:
			def, ok := skillRegistry.Get(item.Name)
			if !ok {
				def = capability.SkillDefinition{
					Name:        item.Name,
					Description: "分析简历文本，提取候选人的优势、风险和优化建议",
					Enabled:     true,
					Kind:        "skill",
					Tags:        []string{"resume", "analysis", "skill"},
				}
			}
			registry.MustRegister(skills.NewResumeAnalyzerSkill(def))
		}
	}

	for _, item := range capCfg.Tools {
		if !item.Enabled {
			continue
		}

		switch item.Name {
		case capability.CapabilityKeywordExtract:
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

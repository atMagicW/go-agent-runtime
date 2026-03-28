package app

import (
	"strings"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

// CapabilityService 提供能力列表查询
type CapabilityService struct {
	registry capabilityRegistryWithList
}

type capabilityRegistryWithList interface {
	ListDescriptors() []capability.Descriptor
}

// NewCapabilityService 创建 CapabilityService
func NewCapabilityService(registry capabilityRegistryWithList) *CapabilityService {
	return &CapabilityService{
		registry: registry,
	}
}

// ListCapabilities 列出所有已注册能力
func (s *CapabilityService) ListCapabilities() []httpapi.CapabilityView {
	if s.registry == nil {
		return nil
	}

	descriptors := s.registry.ListDescriptors()
	out := make([]httpapi.CapabilityView, 0, len(descriptors))
	for _, item := range descriptors {
		view := httpapi.CapabilityView{
			Name:        item.Name,
			Kind:        string(item.Kind),
			Description: item.Description,
			Tags:        item.Tags,
			Version:     item.Version,
			Enabled:     item.Enabled,
			Source:      capabilitySource(item),
		}

		if item.Kind == capability.KindMCPTool {
			view.ServerName = capabilityServerName(item.Tags)
			view.RemoteTool = capabilityRemoteTool(item.Name)
		}

		out = append(out, view)
	}

	return out
}

func capabilitySource(item capability.Descriptor) string {
	switch item.Kind {
	case capability.KindSkill, capability.KindTool:
		return "local"
	case capability.KindMCPTool:
		return "remote"
	default:
		return "unknown"
	}
}

func capabilityServerName(tags []string) string {
	for _, tag := range tags {
		if tag == "" || tag == "mcp" || tag == "remote_tool" {
			continue
		}
		return tag
	}
	return ""
}

func capabilityRemoteTool(name string) string {
	if strings.HasPrefix(name, "mcp_") {
		return strings.TrimPrefix(name, "mcp_")
	}
	return ""
}

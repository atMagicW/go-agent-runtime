package app

import "github.com/atMagicW/go-agent-runtime/internal/domain/capability"

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
func (s *CapabilityService) ListCapabilities() []capability.Descriptor {
	if s.registry == nil {
		return nil
	}
	return s.registry.ListDescriptors()
}

package router

import (
	"context"
	"fmt"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// CapabilityRouter 是 Skill / Tool / MCP 的统一路由器
type CapabilityRouter struct {
	registry portsCapabilityRegistry
}

// portsCapabilityRegistry 是本文件内最小依赖抽象
type portsCapabilityRegistry interface {
	Get(name string) (ports.Capability, bool)
}

// NewCapabilityRouter 创建能力路由器
func NewCapabilityRouter(registry portsCapabilityRegistry) *CapabilityRouter {
	return &CapabilityRouter{
		registry: registry,
	}
}

// Invoke 调用能力
func (r *CapabilityRouter) Invoke(
	ctx context.Context,
	_ agent.RuntimeContext,
	req ports.CapabilityCallRequest,
) (ports.CapabilityCallResponse, error) {
	if r.registry == nil {
		return ports.CapabilityCallResponse{}, fmt.Errorf("capability registry is nil")
	}

	capabilityImpl, ok := r.registry.Get(req.Name)
	if !ok {
		return ports.CapabilityCallResponse{}, fmt.Errorf("capability not found: %s", req.Name)
	}

	result, err := capabilityImpl.Invoke(ctx, req.Input)
	if err != nil {
		return ports.CapabilityCallResponse{}, err
	}

	if !result.Success {
		return ports.CapabilityCallResponse{}, fmt.Errorf("capability failed: %s", result.Error)
	}

	return ports.CapabilityCallResponse{
		Output: result.Output,
	}, nil
}

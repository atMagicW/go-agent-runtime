package router

import (
	"context"
	"fmt"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// CapabilityRouter 是第一版 Skill / MCP 路由器
type CapabilityRouter struct {
}

// NewCapabilityRouter 创建能力路由器
func NewCapabilityRouter() *CapabilityRouter {
	return &CapabilityRouter{}
}

// Invoke 调用能力
func (r *CapabilityRouter) Invoke(
	_ context.Context,
	_ agent.RuntimeContext,
	req ports.CapabilityCallRequest,
) (ports.CapabilityCallResponse, error) {
	return ports.CapabilityCallResponse{
		Output: map[string]any{
			"capability_name": req.Name,
			"result":          fmt.Sprintf("capability %s executed", req.Name),
		},
	}, nil
}

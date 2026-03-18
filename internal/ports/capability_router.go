package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// CapabilityCallRequest 表示 Skill / Tool / MCP 调用请求
type CapabilityCallRequest struct {
	Name  string
	Input map[string]any
}

// CapabilityCallResponse 表示能力调用响应
type CapabilityCallResponse struct {
	Output map[string]any
}

// CapabilityRouter 定义能力路由接口
type CapabilityRouter interface {
	Invoke(ctx context.Context, runtimeCtx agent.RuntimeContext, req CapabilityCallRequest) (CapabilityCallResponse, error)
}

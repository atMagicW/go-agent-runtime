package router

import (
	"context"
	"fmt"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// ModelRouter 是第一版模型路由器
type ModelRouter struct {
}

// NewModelRouter 创建模型路由器
func NewModelRouter() *ModelRouter {
	return &ModelRouter{}
}

// Generate 执行模型调用
func (r *ModelRouter) Generate(
	_ context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ModelCallRequest,
) (ports.ModelCallResponse, error) {
	selectedModel := r.selectModel(runtimeCtx, req)

	// 第一版先返回 mock 数据，后面再接 OpenAI / Anthropic / Gemini
	text := fmt.Sprintf("model=%s response for task=%s", selectedModel, req.TaskType)

	return ports.ModelCallResponse{
		Text:     text,
		Tokens:   128,
		Cost:     0.0021,
		Model:    selectedModel,
		Provider: "mock",
	}, nil
}

// selectModel 选择模型
func (r *ModelRouter) selectModel(runtimeCtx agent.RuntimeContext, req ports.ModelCallRequest) string {
	// 1. 用户强制指定模型优先
	if runtimeCtx.Request.Model != "" {
		return runtimeCtx.Request.Model
	}

	// 2. 按任务类型路由
	switch req.TaskType {
	case "intent":
		return "gpt-4.1-mini"
	case "retrieve_answer":
		return "gpt-4.1"
	case "analysis":
		return "gpt-4.1"
	case "write":
		return "claude-sonnet"
	default:
		return "gpt-4.1-mini"
	}
}

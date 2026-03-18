package router

import (
	"context"
	"fmt"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// ModelRouter 是模型路由器
type ModelRouter struct {
	clients map[string]ports.LLMClient
}

// NewModelRouter 创建模型路由器
func NewModelRouter(clients map[string]ports.LLMClient) *ModelRouter {
	return &ModelRouter{
		clients: clients,
	}
}

// Generate 执行模型调用
func (r *ModelRouter) Generate(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ModelCallRequest,
) (ports.ModelCallResponse, error) {
	selectedModel := r.selectModel(runtimeCtx, req)
	selectedProvider := r.selectProvider(selectedModel)

	client, ok := r.clients[selectedProvider]
	if !ok {
		return ports.ModelCallResponse{}, fmt.Errorf("llm provider client not found: %s", selectedProvider)
	}

	resp, err := client.Generate(ctx, ports.LLMGenerateRequest{
		Model:  selectedModel,
		Prompt: req.Prompt,
		Stream: req.Stream,
	})
	if err != nil {
		return ports.ModelCallResponse{}, err
	}

	return ports.ModelCallResponse{
		Text:     resp.Text,
		Tokens:   resp.TotalTokens,
		Cost:     resp.Cost,
		Model:    resp.Model,
		Provider: resp.Provider,
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
		return "gpt-4.1"
	case "llm_generate":
		return "gpt-4.1-mini"
	case "llm_analyze":
		return "gpt-4.1"
	default:
		return "gpt-4.1-mini"
	}
}

// selectProvider 根据模型名选择 provider
func (r *ModelRouter) selectProvider(model string) string {
	switch {
	case model == "":
		return "openai"
	case len(model) >= 3 && model[:3] == "gpt":
		return "openai"
	case len(model) >= 6 && model[:6] == "claude":
		return "anthropic"
	case len(model) >= 6 && model[:6] == "gemini":
		return "gemini"
	default:
		return "openai"
	}
}

package router

import (
	"context"
	"fmt"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	"github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
)

// ModelRouter 是模型路由器
type ModelRouter struct {
	clients   map[string]ports.LLMClient
	breakers  *governance.BreakerRegistry
	fallbacks *governance.FallbackPolicy
}

// NewModelRouter 创建模型路由器
func NewModelRouter(
	clients map[string]ports.LLMClient,
	breakers *governance.BreakerRegistry,
	fallbacks *governance.FallbackPolicy,
) *ModelRouter {
	return &ModelRouter{
		clients:   clients,
		breakers:  breakers,
		fallbacks: fallbacks,
	}
}

// Generate 执行模型调用
func (r *ModelRouter) Generate(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ModelCallRequest,
) (ports.ModelCallResponse, error) {
	primaryModel := r.selectModel(runtimeCtx, req)

	// 按“主模型 + fallback 模型链”依次尝试
	candidates := []string{primaryModel}
	candidates = append(candidates, r.fallbacks.NextModels(primaryModel)...)

	var lastErr error

	for _, modelName := range candidates {
		provider := r.selectProvider(modelName)
		client, ok := r.clients[provider]
		if !ok {
			lastErr = fmt.Errorf("llm provider client not found: %s", provider)
			continue
		}

		breakerName := "model:" + provider + ":" + modelName
		breaker := r.breakers.GetOrCreate(breakerName, 3, 10*time.Second)

		if !breaker.Allow() {
			lastErr = fmt.Errorf("circuit breaker open for model %s", modelName)
			continue
		}

		resp, err := client.Generate(ctx, ports.LLMGenerateRequest{
			Model:  modelName,
			Prompt: req.Prompt,
			Stream: false,
		})
		if err != nil {
			breaker.OnFailure()
			lastErr = err
			continue
		}

		breaker.OnSuccess()

		return ports.ModelCallResponse{
			Text:     resp.Text,
			Tokens:   resp.TotalTokens,
			Cost:     resp.Cost,
			Model:    resp.Model,
			Provider: resp.Provider,
		}, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all model candidates failed")
	}
	return ports.ModelCallResponse{}, lastErr
}

// GenerateStream 流式执行模型调用
func (r *ModelRouter) GenerateStream(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ModelCallRequest,
	onToken ports.ModelStreamHandler,
) error {
	primaryModel := r.selectModel(runtimeCtx, req)
	candidates := []string{primaryModel}
	candidates = append(candidates, r.fallbacks.NextModels(primaryModel)...)

	var lastErr error

	for _, modelName := range candidates {
		provider := r.selectProvider(modelName)
		client, ok := r.clients[provider]
		if !ok {
			lastErr = fmt.Errorf("llm provider client not found: %s", provider)
			continue
		}

		streamingClient, ok := client.(ports.StreamingLLMClient)
		if !ok {
			lastErr = fmt.Errorf("provider %s does not support streaming", provider)
			continue
		}

		breakerName := "model:" + provider + ":" + modelName
		breaker := r.breakers.GetOrCreate(breakerName, 3, 10*time.Second)

		if !breaker.Allow() {
			lastErr = fmt.Errorf("circuit breaker open for model %s", modelName)
			continue
		}

		err := streamingClient.GenerateStream(ctx, ports.LLMGenerateRequest{
			Model:  modelName,
			Prompt: req.Prompt,
			Stream: true,
		}, func(chunk ports.StreamChunk) error {
			if chunk.Done {
				return nil
			}
			return onToken(chunk.Text)
		})
		if err != nil {
			breaker.OnFailure()
			lastErr = err
			continue
		}

		breaker.OnSuccess()
		return nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all streaming model candidates failed")
	}
	return lastErr
}

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

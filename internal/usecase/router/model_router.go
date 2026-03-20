package router

import (
	"context"
	"fmt"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	domainmodel "github.com/atMagicW/go-agent-runtime/internal/domain/model"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	"github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
)

// modelRegistry 是本文件依赖的最小模型注册接口
type modelRegistry interface {
	DefaultModel() string
	Get(name string) (domainmodel.Profile, bool)
	IsEnabled(name string) bool
	ProviderOf(name string) string
	ResolveByTaskType(taskType string) (domainmodel.Profile, bool)
}

// ModelRouter 是模型路由器
type ModelRouter struct {
	clients   map[string]ports.LLMClient
	registry  modelRegistry
	breakers  *governance.BreakerRegistry
	fallbacks *governance.FallbackPolicy
}

// NewModelRouter 创建模型路由器
func NewModelRouter(
	clients map[string]ports.LLMClient,
	registry modelRegistry,
	breakers *governance.BreakerRegistry,
	fallbacks *governance.FallbackPolicy,
) *ModelRouter {
	return &ModelRouter{
		clients:   clients,
		registry:  registry,
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
	primaryModel, err := r.selectModel(runtimeCtx, req)
	if err != nil {
		return ports.ModelCallResponse{}, err
	}

	candidates := []string{primaryModel}
	candidates = append(candidates, r.filterEnabledFallbackModels(primaryModel)...)

	var lastErr error

	for _, modelName := range candidates {
		provider := r.selectProvider(modelName)
		if provider == "" {
			lastErr = fmt.Errorf("provider not found for model %s", modelName)
			continue
		}

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
	primaryModel, err := r.selectModel(runtimeCtx, req)
	if err != nil {
		return err
	}

	candidates := []string{primaryModel}
	candidates = append(candidates, r.filterEnabledFallbackModels(primaryModel)...)

	var lastErr error

	for _, modelName := range candidates {
		provider := r.selectProvider(modelName)
		if provider == "" {
			lastErr = fmt.Errorf("provider not found for model %s", modelName)
			continue
		}

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

func (r *ModelRouter) selectModel(
	runtimeCtx agent.RuntimeContext,
	req ports.ModelCallRequest,
) (string, error) {
	if runtimeCtx.Request.Model != "" {
		if r.registry != nil && r.registry.IsEnabled(runtimeCtx.Request.Model) {
			return runtimeCtx.Request.Model, nil
		}
		return "", fmt.Errorf("requested model is not enabled: %s", runtimeCtx.Request.Model)
	}

	if r.registry != nil {
		if item, ok := r.registry.ResolveByTaskType(req.TaskType); ok {
			return item.Name, nil
		}
	}

	if r.registry != nil {
		model := r.registry.DefaultModel()
		if model != "" && r.registry.IsEnabled(model) {
			return model, nil
		}
	}

	return "", fmt.Errorf("no enabled model available")
}

func (r *ModelRouter) selectProvider(model string) string {
	if r.registry != nil {
		if provider := r.registry.ProviderOf(model); provider != "" {
			return provider
		}
	}
	return ""
}

func (r *ModelRouter) filterEnabledFallbackModels(primaryModel string) []string {
	raw := r.fallbacks.NextModels(primaryModel)
	if r.registry == nil {
		return raw
	}

	out := make([]string, 0, len(raw))
	seen := map[string]struct{}{}

	for _, name := range raw {
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		if r.registry.IsEnabled(name) {
			out = append(out, name)
			seen[name] = struct{}{}
		}
	}

	return out
}

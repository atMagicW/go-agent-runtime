package router

import (
	"context"
	"fmt"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	"github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
)

// portsCapabilityRegistry 是本文件内最小依赖抽象
type portsCapabilityRegistry interface {
	Get(name string) (ports.Capability, bool)
}

// CapabilityRouter 是 Skill / Tool / MCP 的统一路由器
type CapabilityRouter struct {
	registry  portsCapabilityRegistry
	breakers  *governance.BreakerRegistry
	fallbacks *governance.FallbackPolicy
}

// NewCapabilityRouter 创建能力路由器
func NewCapabilityRouter(
	registry portsCapabilityRegistry,
	breakers *governance.BreakerRegistry,
	fallbacks *governance.FallbackPolicy,
) *CapabilityRouter {
	return &CapabilityRouter{
		registry:  registry,
		breakers:  breakers,
		fallbacks: fallbacks,
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

	candidates := []string{req.Name}
	candidates = append(candidates, r.fallbacks.NextCapabilities(req.Name)...)

	var lastErr error

	for _, capabilityName := range candidates {
		capabilityImpl, ok := r.registry.Get(capabilityName)
		if !ok {
			lastErr = fmt.Errorf("capability not found: %s", capabilityName)
			continue
		}

		breaker := r.breakers.GetOrCreate("capability:"+capabilityName, 3, 10*time.Second)
		if !breaker.Allow() {
			lastErr = fmt.Errorf("circuit breaker open for capability %s", capabilityName)
			continue
		}

		input := cloneMap(req.Input)
		input["name"] = capabilityName

		result, err := capabilityImpl.Invoke(ctx, input)
		if err != nil {
			breaker.OnFailure()
			lastErr = err
			continue
		}

		if !result.Success {
			breaker.OnFailure()
			lastErr = fmt.Errorf("capability failed: %s", result.Error)
			continue
		}

		breaker.OnSuccess()

		return ports.CapabilityCallResponse{
			Output: result.Output,
		}, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all capability candidates failed")
	}
	return ports.CapabilityCallResponse{}, lastErr
}

func cloneMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

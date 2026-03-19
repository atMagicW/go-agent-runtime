package router

import (
	"context"
	"fmt"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	"github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
)

// RAGRouter 是检索路由器
type RAGRouter struct {
	breakers  *governance.BreakerRegistry
	fallbacks *governance.FallbackPolicy
}

// NewRAGRouter 创建检索路由器
func NewRAGRouter(
	breakers *governance.BreakerRegistry,
	fallbacks *governance.FallbackPolicy,
) *RAGRouter {
	return &RAGRouter{
		breakers:  breakers,
		fallbacks: fallbacks,
	}
}

// Retrieve 执行检索
func (r *RAGRouter) Retrieve(
	_ context.Context,
	_ agent.RuntimeContext,
	req ports.RetrievalRequest,
) (ports.RetrievalResponse, error) {
	candidates := []string{req.KnowledgeBase}
	candidates = append(candidates, r.fallbacks.NextKnowledgeBases(req.KnowledgeBase)...)

	var lastErr error

	for _, kb := range candidates {
		breaker := r.breakers.GetOrCreate("rag:"+kb, 3, 10*time.Second)
		if !breaker.Allow() {
			lastErr = fmt.Errorf("circuit breaker open for knowledge base %s", kb)
			continue
		}

		// 第一版先保留 mock 检索逻辑
		resp := ports.RetrievalResponse{
			Evidences: []map[string]any{
				{
					"kb":      kb,
					"content": fmt.Sprintf("evidence-1 for query: %s", req.Query),
					"score":   0.91,
				},
				{
					"kb":      kb,
					"content": fmt.Sprintf("evidence-2 for query: %s", req.Query),
					"score":   0.84,
				},
			},
		}

		breaker.OnSuccess()
		return resp, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all knowledge base candidates failed")
	}
	return ports.RetrievalResponse{}, lastErr
}

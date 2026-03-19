package router

import (
	"context"
	"fmt"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	"github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
)

// ragSearchService 是本文件内最小依赖接口
type ragSearchService interface {
	Search(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error)
}

// RAGRouter 是检索路由器
type RAGRouter struct {
	service   ragSearchService
	breakers  *governance.BreakerRegistry
	fallbacks *governance.FallbackPolicy
}

// NewRAGRouter 创建检索路由器
func NewRAGRouter(
	service ragSearchService,
	breakers *governance.BreakerRegistry,
	fallbacks *governance.FallbackPolicy,
) *RAGRouter {
	return &RAGRouter{
		service:   service,
		breakers:  breakers,
		fallbacks: fallbacks,
	}
}

// Retrieve 执行检索
func (r *RAGRouter) Retrieve(
	ctx context.Context,
	_ agent.RuntimeContext,
	req ports.RetrievalRequest,
) (ports.RetrievalResponse, error) {
	if r.service == nil {
		return ports.RetrievalResponse{}, fmt.Errorf("rag service is nil")
	}

	candidates := []string{req.KnowledgeBase}
	candidates = append(candidates, r.fallbacks.NextKnowledgeBases(req.KnowledgeBase)...)

	var lastErr error

	for _, kb := range candidates {
		breaker := r.breakers.GetOrCreate("rag:"+kb, 3, 10*time.Second)
		if !breaker.Allow() {
			lastErr = fmt.Errorf("circuit breaker open for knowledge base %s", kb)
			continue
		}

		items, err := r.service.Search(ctx, kb, req.Query, req.TopK)
		if err != nil {
			breaker.OnFailure()
			lastErr = err
			continue
		}

		breaker.OnSuccess()

		evidences := make([]map[string]any, 0, len(items))
		for _, item := range items {
			evidences = append(evidences, map[string]any{
				"kb":       item.KBID,
				"doc_id":   item.DocID,
				"chunk_id": item.ChunkID,
				"content":  item.Content,
				"score":    item.Score,
				"metadata": item.Metadata,
			})
		}

		return ports.RetrievalResponse{
			Evidences: evidences,
		}, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all knowledge base candidates failed")
	}
	return ports.RetrievalResponse{}, lastErr
}

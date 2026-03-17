package router

import (
	"context"
	"fmt"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// RAGRouter 是第一版检索路由器
type RAGRouter struct {
}

// NewRAGRouter 创建检索路由器
func NewRAGRouter() *RAGRouter {
	return &RAGRouter{}
}

// Retrieve 执行检索
func (r *RAGRouter) Retrieve(
	_ context.Context,
	_ agent.RuntimeContext,
	req ports.RetrievalRequest,
) (ports.RetrievalResponse, error) {
	return ports.RetrievalResponse{
		Evidences: []map[string]any{
			{
				"kb":      req.KnowledgeBase,
				"content": fmt.Sprintf("evidence-1 for query: %s", req.Query),
				"score":   0.91,
			},
			{
				"kb":      req.KnowledgeBase,
				"content": fmt.Sprintf("evidence-2 for query: %s", req.Query),
				"score":   0.84,
			},
		},
	}, nil
}

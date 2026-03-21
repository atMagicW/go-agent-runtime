package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

// Reranker 定义检索结果重排接口
type Reranker interface {
	Rerank(ctx context.Context, query string, items []rag.Evidence, topK int) ([]rag.Evidence, error)
}

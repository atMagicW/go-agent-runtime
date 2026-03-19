package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

// RAGRepository 定义 RAG 持久化与检索接口
type RAGRepository interface {
	EnsureKnowledgeBase(ctx context.Context, kb rag.KnowledgeBase) error

	InsertDocument(ctx context.Context, doc rag.Document) error

	InsertChunk(ctx context.Context, chunk rag.Chunk) error

	SearchByVector(ctx context.Context, kbID string, embedding []float32, topK int) ([]rag.Evidence, error)

	SearchByKeyword(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error)
}

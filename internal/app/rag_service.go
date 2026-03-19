package app

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// RAGService 提供真实检索能力
type RAGService struct {
	repo      ports.RAGRepository
	embedding ports.EmbeddingProvider
}

// NewRAGService 创建 RAGService
func NewRAGService(repo ports.RAGRepository, embedding ports.EmbeddingProvider) *RAGService {
	return &RAGService{
		repo:      repo,
		embedding: embedding,
	}
}

// Search 执行向量检索，失败时回退关键词检索
func (s *RAGService) Search(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error) {
	vector, err := s.embedding.Embed(ctx, query)
	if err == nil && len(vector) > 0 {
		items, searchErr := s.repo.SearchByVector(ctx, kbID, vector, topK)
		if searchErr == nil && len(items) > 0 {
			return items, nil
		}
	}

	return s.repo.SearchByKeyword(ctx, kbID, query, topK)
}

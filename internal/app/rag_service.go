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
	reranker  ports.Reranker
}

// NewRAGService 创建 RAGService
func NewRAGService(
	repo ports.RAGRepository,
	embedding ports.EmbeddingProvider,
	reranker ports.Reranker,
) *RAGService {
	return &RAGService{
		repo:      repo,
		embedding: embedding,
		reranker:  reranker,
	}
}

// Search 执行二阶段检索：召回 + 重排
func (s *RAGService) Search(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error) {
	if topK <= 0 {
		topK = 5
	}

	candidates := make([]rag.Evidence, 0)

	// 1. 向量召回
	vector, err := s.embedding.Embed(ctx, query)
	if err == nil && len(vector) > 0 {
		vectorItems, searchErr := s.repo.SearchByVector(ctx, kbID, vector, topK*2)
		if searchErr == nil {
			candidates = append(candidates, vectorItems...)
		}
	}

	// 2. 关键词召回
	keywordItems, keywordErr := s.repo.SearchByKeyword(ctx, kbID, query, topK*2)
	if keywordErr == nil {
		candidates = append(candidates, keywordItems...)
	}

	// 3. 去重
	candidates = dedupEvidence(candidates)

	// 4. 没有 reranker 时直接截断
	if s.reranker == nil {
		if len(candidates) > topK {
			return candidates[:topK], nil
		}
		return candidates, nil
	}

	// 5. rerank
	reranked, err := s.reranker.Rerank(ctx, query, candidates, topK)
	if err != nil {
		if len(candidates) > topK {
			return candidates[:topK], nil
		}
		return candidates, nil
	}

	return reranked, nil
}

func dedupEvidence(items []rag.Evidence) []rag.Evidence {
	seen := make(map[string]struct{}, len(items))
	out := make([]rag.Evidence, 0, len(items))

	for _, item := range items {
		key := item.KBID + "::" + item.DocID + "::" + item.ChunkID
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item)
	}

	return out
}

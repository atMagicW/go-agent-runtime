package app

import (
	"context"
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

type mockEmbeddingProvider struct{}

func (m *mockEmbeddingProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	_ = ctx
	_ = text
	return []float32{0.1, 0.2, 0.3}, nil
}

type mockRAGRepository struct {
	vectorItems  []rag.Evidence
	keywordItems []rag.Evidence
}

func (m *mockRAGRepository) EnsureKnowledgeBase(ctx context.Context, kb rag.KnowledgeBase) error {
	_ = ctx
	_ = kb
	return nil
}

func (m *mockRAGRepository) InsertDocument(ctx context.Context, doc rag.Document) error {
	_ = ctx
	_ = doc
	return nil
}

func (m *mockRAGRepository) InsertChunk(ctx context.Context, chunk rag.Chunk) error {
	_ = ctx
	_ = chunk
	return nil
}

func (m *mockRAGRepository) SearchByVector(ctx context.Context, kbID string, embedding []float32, topK int) ([]rag.Evidence, error) {
	_ = ctx
	_ = kbID
	_ = embedding
	_ = topK
	return m.vectorItems, nil
}

func (m *mockRAGRepository) SearchByKeyword(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error) {
	_ = ctx
	_ = kbID
	_ = query
	_ = topK
	return m.keywordItems, nil
}

type mockReranker struct{}

func (m *mockReranker) Rerank(ctx context.Context, query string, items []rag.Evidence, topK int) ([]rag.Evidence, error) {
	_ = ctx
	_ = query
	if len(items) > topK {
		return items[:topK], nil
	}
	return items, nil
}

func TestRAGService_Search(t *testing.T) {
	repo := &mockRAGRepository{
		vectorItems: []rag.Evidence{
			{KBID: "default", DocID: "d1", ChunkID: "c1", Content: "vector result", Score: 0.9},
		},
		keywordItems: []rag.Evidence{
			{KBID: "default", DocID: "d2", ChunkID: "c2", Content: "keyword result", Score: 0.8},
		},
	}

	svc := NewRAGService(repo, &mockEmbeddingProvider{}, &mockReranker{})

	out, err := svc.Search(context.Background(), "default", "agent runtime", 2)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(out) != 2 {
		t.Fatalf("Search() len = %d, want 2", len(out))
	}
}

func TestDedupEvidence(t *testing.T) {
	items := []rag.Evidence{
		{KBID: "default", DocID: "d1", ChunkID: "c1"},
		{KBID: "default", DocID: "d1", ChunkID: "c1"},
		{KBID: "default", DocID: "d1", ChunkID: "c2"},
	}

	out := dedupEvidence(items)
	if len(out) != 2 {
		t.Fatalf("dedupEvidence() len = %d, want 2", len(out))
	}
	t.Log(out)
}

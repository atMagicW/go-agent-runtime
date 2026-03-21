package app

import (
	"context"
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/textsplitter"
)

type mockIngestRepo struct {
	docCount   int
	chunkCount int
}

func (m *mockIngestRepo) EnsureKnowledgeBase(ctx context.Context, kb rag.KnowledgeBase) error {
	_ = ctx
	_ = kb
	return nil
}

func (m *mockIngestRepo) InsertDocument(ctx context.Context, doc rag.Document) error {
	_ = ctx
	_ = doc
	m.docCount++
	return nil
}

func (m *mockIngestRepo) InsertChunk(ctx context.Context, chunk rag.Chunk) error {
	_ = ctx
	_ = chunk
	m.chunkCount++
	return nil
}

func (m *mockIngestRepo) SearchByVector(ctx context.Context, kbID string, embedding []float32, topK int) ([]rag.Evidence, error) {
	_ = ctx
	_ = kbID
	_ = embedding
	_ = topK
	return nil, nil
}

func (m *mockIngestRepo) SearchByKeyword(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error) {
	_ = ctx
	_ = kbID
	_ = query
	_ = topK
	return nil, nil
}

func TestIngestText(t *testing.T) {
	repo := &mockIngestRepo{}
	embedding := &mockEmbeddingProvider{}
	splitter := textsplitter.NewSplitter(10, 2)

	svc := NewIngestService(repo, embedding, splitter)

	resp, err := svc.IngestText(context.Background(), rag.IngestTextRequest{
		KBID:  "default",
		Title: "test doc",
		Text:  "abcdefghijklmnopqrstuvwxyz",
	})
	if err != nil {
		t.Fatalf("IngestText() error = %v", err)
	}

	if resp.DocID == "" {
		t.Fatal("DocID is empty")
	}

	if repo.docCount != 1 {
		t.Fatalf("docCount = %d, want 1", repo.docCount)
	}

	if repo.chunkCount < 2 {
		t.Fatalf("chunkCount = %d, want >= 2", repo.chunkCount)
	}
	t.Log(repo)
}

package rerank

import (
	"context"
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

func TestSimpleReranker_Rerank(t *testing.T) {
	r := NewSimpleReranker()

	items := []rag.Evidence{
		{
			KBID:    "default",
			DocID:   "d1",
			ChunkID: "c1",
			Content: "golang agent runtime architecture",
			Score:   0.6,
		},
		{
			KBID:    "default",
			DocID:   "d1",
			ChunkID: "c2",
			Content: "cooking recipe and food notes",
			Score:   0.9,
		},
	}

	out, err := r.Rerank(context.Background(), "golang agent runtime", items, 2)
	if err != nil {
		t.Fatalf("Rerank error = %v", err)
	}

	if len(out) != 2 {
		t.Fatalf("len(out) = %d, want 2", len(out))
	}

	if out[0].ChunkID != "c1" {
		t.Fatalf("top chunk = %s, want c1", out[0].ChunkID)
	}
}

package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

type Repository struct {
	mu     sync.RWMutex
	kbs    map[string]rag.KnowledgeBase
	docs   map[string]rag.Document
	chunks []rag.Chunk
}

func NewRepository() *Repository {
	return &Repository{
		kbs:    map[string]rag.KnowledgeBase{},
		docs:   map[string]rag.Document{},
		chunks: []rag.Chunk{},
	}
}

func (r *Repository) EnsureKnowledgeBase(ctx context.Context, kb rag.KnowledgeBase) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.kbs[kb.KBID] = kb
	return nil
}

func (r *Repository) InsertDocument(ctx context.Context, doc rag.Document) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.docs[doc.DocID] = doc
	return nil
}

func (r *Repository) InsertChunk(ctx context.Context, chunk rag.Chunk) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.chunks = append(r.chunks, chunk)
	return nil
}

func (r *Repository) SearchByVector(ctx context.Context, kbID string, embedding []float32, topK int) ([]rag.Evidence, error) {
	_ = ctx
	_ = embedding

	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]rag.Evidence, 0)
	for _, c := range r.chunks {
		if c.KBID != kbID {
			continue
		}
		out = append(out, rag.Evidence{
			KBID:     c.KBID,
			DocID:    c.DocID,
			ChunkID:  c.ChunkID,
			Content:  c.Content,
			Score:    0.8,
			Metadata: c.Metadata,
		})
		if topK > 0 && len(out) >= topK {
			break
		}
	}
	return out, nil
}

func (r *Repository) SearchByKeyword(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	lq := strings.ToLower(query)
	out := make([]rag.Evidence, 0)

	for _, c := range r.chunks {
		if c.KBID != kbID {
			continue
		}
		if !strings.Contains(strings.ToLower(c.Content), lq) {
			continue
		}

		out = append(out, rag.Evidence{
			KBID:     c.KBID,
			DocID:    c.DocID,
			ChunkID:  c.ChunkID,
			Content:  c.Content,
			Score:    0.9,
			Metadata: c.Metadata,
		})

		if topK > 0 && len(out) >= topK {
			break
		}
	}

	return out, nil
}

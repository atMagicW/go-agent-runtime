package pgvector

import (
	"context"

	"github.com/google/uuid"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// SeedDemoKnowledgeBase 写入演示数据
func SeedDemoKnowledgeBase(
	ctx context.Context,
	repo *Repository,
	embedding ports.EmbeddingProvider,
	kbID string,
) error {
	docID := uuid.New().String()

	doc := rag.Document{
		DocID:  docID,
		KBID:   kbID,
		Title:  "Go Agent Runtime Intro",
		Source: "demo",
		Metadata: map[string]any{
			"category": "architecture",
		},
	}

	if err := repo.InsertDocument(ctx, doc); err != nil {
		return err
	}

	chunks := []string{
		"Go agent runtime supports intent recognition, planning, model routing and capability dispatch.",
		"RAG routing can select different knowledge bases based on intent and context relevance.",
		"Capability registry can unify local skill, local tool and remote MCP tool under one runtime.",
	}

	for _, content := range chunks {
		vec, err := embedding.Embed(ctx, content)
		if err != nil {
			return err
		}

		err = repo.InsertChunk(ctx, rag.Chunk{
			ChunkID:   uuid.New().String(),
			DocID:     docID,
			KBID:      kbID,
			Content:   content,
			Metadata:  map[string]any{"seed": true},
			Embedding: vec,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

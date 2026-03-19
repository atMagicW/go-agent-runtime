package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	"github.com/atMagicW/go-agent-runtime/internal/pkg/textsplitter"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// IngestService 提供知识库写入能力
type IngestService struct {
	repo      ports.RAGRepository
	embedding ports.EmbeddingProvider
	splitter  *textsplitter.Splitter
}

// NewIngestService 创建 IngestService
func NewIngestService(
	repo ports.RAGRepository,
	embedding ports.EmbeddingProvider,
	splitter *textsplitter.Splitter,
) *IngestService {
	return &IngestService{
		repo:      repo,
		embedding: embedding,
		splitter:  splitter,
	}
}

// IngestText 将一段文本切块后写入知识库
func (s *IngestService) IngestText(ctx context.Context, req rag.IngestTextRequest) (*rag.IngestTextResponse, error) {
	if req.KBID == "" {
		return nil, fmt.Errorf("kb_id is required")
	}
	if req.Text == "" {
		return nil, fmt.Errorf("text is required")
	}
	if req.Title == "" {
		req.Title = "Untitled Document"
	}
	if req.TenantID == "" {
		req.TenantID = "default"
	}

	err := s.repo.EnsureKnowledgeBase(ctx, rag.KnowledgeBase{
		KBID:        req.KBID,
		TenantID:    req.TenantID,
		Name:        req.KBID,
		Description: req.Description,
		Enabled:     true,
	})
	if err != nil {
		return nil, err
	}

	docID := uuid.New().String()
	doc := rag.Document{
		DocID:    docID,
		KBID:     req.KBID,
		Title:    req.Title,
		Source:   req.Source,
		Metadata: req.Metadata,
	}

	if err := s.repo.InsertDocument(ctx, doc); err != nil {
		return nil, err
	}

	chunks := s.splitter.Split(req.Text)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunk generated")
	}

	for idx, content := range chunks {
		vec, err := s.embedding.Embed(ctx, content)
		if err != nil {
			return nil, err
		}

		err = s.repo.InsertChunk(ctx, rag.Chunk{
			ChunkID:   uuid.New().String(),
			DocID:     docID,
			KBID:      req.KBID,
			Content:   content,
			Embedding: vec,
			Metadata: map[string]any{
				"chunk_index": idx,
				"title":       req.Title,
				"source":      req.Source,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return &rag.IngestTextResponse{
		KBID:       req.KBID,
		DocID:      docID,
		ChunkCount: len(chunks),
	}, nil
}

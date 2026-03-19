package pgvector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

// Repository 是基于 PostgreSQL / pgvector 的 RAG 仓储实现
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository 创建 pgvector RAG 仓储
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// EnsureKnowledgeBase 确保知识库存在
func (r *Repository) EnsureKnowledgeBase(ctx context.Context, kb rag.KnowledgeBase) error {
	const query = `
INSERT INTO knowledge_bases (kb_id, tenant_id, name, description, enabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (kb_id)
DO UPDATE SET
    tenant_id = EXCLUDED.tenant_id,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    enabled = EXCLUDED.enabled,
    updated_at = NOW()
`
	_, err := r.db.Exec(ctx, query, kb.KBID, kb.TenantID, kb.Name, kb.Description, kb.Enabled)
	return err
}

// InsertDocument 插入文档
func (r *Repository) InsertDocument(ctx context.Context, doc rag.Document) error {
	if doc.DocID == "" {
		doc.DocID = uuid.New().String()
	}

	metadataJSON, err := json.Marshal(doc.Metadata)
	if err != nil {
		return err
	}

	const query = `
INSERT INTO kb_documents (doc_id, kb_id, title, source, metadata_json, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
ON CONFLICT (doc_id) DO NOTHING
`
	_, err = r.db.Exec(ctx, query, doc.DocID, doc.KBID, doc.Title, doc.Source, metadataJSON)
	return err
}

// InsertChunk 插入切块
func (r *Repository) InsertChunk(ctx context.Context, chunk rag.Chunk) error {
	if chunk.ChunkID == "" {
		chunk.ChunkID = uuid.New().String()
	}

	metadataJSON, err := json.Marshal(chunk.Metadata)
	if err != nil {
		return err
	}

	vec := pgvector.NewVector(chunk.Embedding)

	const query = `
INSERT INTO kb_document_chunks (chunk_id, doc_id, kb_id, content, metadata_json, embedding, created_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
ON CONFLICT (chunk_id) DO NOTHING
`
	_, err = r.db.Exec(ctx, query, chunk.ChunkID, chunk.DocID, chunk.KBID, chunk.Content, metadataJSON, vec)
	return err
}

// SearchByVector 向量检索
func (r *Repository) SearchByVector(ctx context.Context, kbID string, embedding []float32, topK int) ([]rag.Evidence, error) {
	if topK <= 0 {
		topK = 5
	}

	vec := pgvector.NewVector(embedding)

	const query = `
SELECT
    kb_id,
    doc_id,
    chunk_id,
    content,
    metadata_json,
    1 - (embedding <=> $2) AS score
FROM kb_document_chunks
WHERE kb_id = $1
ORDER BY embedding <=> $2
LIMIT $3
`
	rows, err := r.db.Query(ctx, query, kbID, vec, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]rag.Evidence, 0)
	for rows.Next() {
		var (
			kb          string
			docID       string
			chunkID     string
			content     string
			metadataRaw []byte
			score       float64
		)

		if err := rows.Scan(&kb, &docID, &chunkID, &content, &metadataRaw, &score); err != nil {
			return nil, err
		}

		metadata := map[string]any{}
		if len(metadataRaw) > 0 {
			_ = json.Unmarshal(metadataRaw, &metadata)
		}

		out = append(out, rag.Evidence{
			KBID:     kb,
			DocID:    docID,
			ChunkID:  chunkID,
			Content:  content,
			Score:    score,
			Metadata: metadata,
		})
	}

	return out, rows.Err()
}

// SearchByKeyword 关键词检索
func (r *Repository) SearchByKeyword(ctx context.Context, kbID string, queryText string, topK int) ([]rag.Evidence, error) {
	if topK <= 0 {
		topK = 5
	}

	const query = `
SELECT
    kb_id,
    doc_id,
    chunk_id,
    content,
    metadata_json
FROM kb_document_chunks
WHERE kb_id = $1
  AND content ILIKE '%' || $2 || '%'
LIMIT $3
`
	rows, err := r.db.Query(ctx, query, kbID, queryText, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]rag.Evidence, 0)
	for rows.Next() {
		var (
			kb          string
			docID       string
			chunkID     string
			content     string
			metadataRaw []byte
		)

		if err := rows.Scan(&kb, &docID, &chunkID, &content, &metadataRaw); err != nil {
			return nil, err
		}

		metadata := map[string]any{}
		if len(metadataRaw) > 0 {
			_ = json.Unmarshal(metadataRaw, &metadata)
		}

		out = append(out, rag.Evidence{
			KBID:     kb,
			DocID:    docID,
			ChunkID:  chunkID,
			Content:  content,
			Score:    0.5,
			Metadata: metadata,
		})
	}

	return out, rows.Err()
}

// SeedDemoData 写入一批演示数据
func (r *Repository) SeedDemoData(ctx context.Context, kbID string, docs []rag.Document, chunks []rag.Chunk) error {
	if kbID == "" {
		return fmt.Errorf("kbID is empty")
	}

	for _, doc := range docs {
		if err := r.InsertDocument(ctx, doc); err != nil {
			return err
		}
	}

	for _, chunk := range chunks {
		if err := r.InsertChunk(ctx, chunk); err != nil {
			return err
		}
	}

	return nil
}

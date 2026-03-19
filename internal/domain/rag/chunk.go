package rag

import "time"

// Chunk 表示文档切块
type Chunk struct {
	ChunkID   string         `json:"chunk_id"`
	DocID     string         `json:"doc_id"`
	KBID      string         `json:"kb_id"`
	Content   string         `json:"content"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Embedding []float32      `json:"embedding,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

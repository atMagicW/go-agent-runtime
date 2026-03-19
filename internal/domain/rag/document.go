package rag

import "time"

// Document 表示知识库中的文档
type Document struct {
	DocID     string         `json:"doc_id"`
	KBID      string         `json:"kb_id"`
	Title     string         `json:"title"`
	Source    string         `json:"source"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

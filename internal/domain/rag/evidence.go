package rag

// Evidence 表示一次检索返回的证据
type Evidence struct {
	KBID     string         `json:"kb_id"`
	DocID    string         `json:"doc_id"`
	ChunkID  string         `json:"chunk_id"`
	Content  string         `json:"content"`
	Score    float64        `json:"score"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

package rag

// IngestTextRequest 表示文本入库请求
type IngestTextRequest struct {
	KBID        string         `json:"kb_id"`
	TenantID    string         `json:"tenant_id"`
	Title       string         `json:"title"`
	Source      string         `json:"source"`
	Text        string         `json:"text"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// IngestTextResponse 表示文本入库结果
type IngestTextResponse struct {
	KBID       string `json:"kb_id"`
	DocID      string `json:"doc_id"`
	ChunkCount int    `json:"chunk_count"`
}

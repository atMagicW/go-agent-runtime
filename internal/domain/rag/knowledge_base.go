package rag

import "time"

// KnowledgeBase 表示一个知识库
type KnowledgeBase struct {
	KBID        string    `json:"kb_id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

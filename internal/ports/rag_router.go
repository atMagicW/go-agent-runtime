package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// RetrievalRequest 表示检索请求
type RetrievalRequest struct {
	KnowledgeBase string
	Query         string
	TopK          int
}

// RetrievalResponse 表示检索结果
type RetrievalResponse struct {
	Evidences []map[string]any
}

// RAGRouter 定义多知识库路由接口
type RAGRouter interface {
	Retrieve(ctx context.Context, runtimeCtx agent.RuntimeContext, req RetrievalRequest) (RetrievalResponse, error)
}

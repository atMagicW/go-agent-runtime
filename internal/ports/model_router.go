package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// ModelCallRequest 表示一次模型调用请求
type ModelCallRequest struct {
	TaskType string
	Model    string
	Prompt   string
	Stream   bool
}

// ModelCallResponse 表示模型调用结果
type ModelCallResponse struct {
	Text   string
	Tokens int
	Cost   float64
}

// ModelRouter 定义多模型路由接口
type ModelRouter interface {
	Generate(ctx context.Context, runtimeCtx agent.RuntimeContext, req ModelCallRequest) (ModelCallResponse, error)
}

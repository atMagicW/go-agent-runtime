package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// ComposeRequest 表示最终回答生成请求
type ComposeRequest struct {
	Message     string
	PromptName  string
	PromptVer   string
	StepResults []agent.StepResult
}

// ComposeResponse 表示最终回答生成结果
type ComposeResponse struct {
	Text   string
	Tokens int
	Cost   float64
	Model  string
}

// ResponseComposer 定义最终回答生成器接口
type ResponseComposer interface {
	Compose(ctx context.Context, runtimeCtx agent.RuntimeContext, req ComposeRequest) (ComposeResponse, error)
}

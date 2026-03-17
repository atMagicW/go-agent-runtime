package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// IntentEngine 定义意图识别能力
type IntentEngine interface {
	Recognize(ctx context.Context, runtimeCtx agent.RuntimeContext, message string) (agent.IntentResult, error)
}

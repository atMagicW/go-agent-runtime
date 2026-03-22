package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// IntentEngine 定义意图识别能力
type IntentEngine interface {
	Recognize(ctx context.Context, runtimeCtx agent.RuntimeContext, message string) (agent.IntentResult, error)
}

// IntentClassifier 定义单个分类器能力
// true：这个分类器已经给出了可用结果
// false：没命中，继续走下一个分类器
type IntentClassifier interface {
	Classify(ctx context.Context, runtimeCtx agent.RuntimeContext, message string) (agent.IntentResult, bool, error)
}

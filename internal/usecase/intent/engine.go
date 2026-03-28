package intent

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// Engine 是组合式意图识别器
type Engine struct {
	classifiers []ports.IntentClassifier
}

// NewEngine 创建默认意图识别器
func NewEngine(classifiers ...ports.IntentClassifier) *Engine {
	if len(classifiers) == 0 {
		classifiers = []ports.IntentClassifier{
			NewRuleClassifier(),
		}
	}

	return &Engine{
		classifiers: classifiers,
	}
}

// Recognize 依次执行多个分类器
func (e *Engine) Recognize(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	message string,
) (agent.IntentResult, error) {
	for _, classifier := range e.classifiers {
		result, ok, err := classifier.Classify(ctx, runtimeCtx, message)
		if err != nil {
			// 单个分类器失败时不立刻中断，继续尝试下一个
			continue
		}
		if ok {
			return result, nil
		}
	}

	// 最终兜底
	return agent.IntentResult{
		IntentType:         agent.IntentChat,
		Confidence:         0.50,
		RequiresRAG:        false,
		RequiresPlanning:   false,
		RequiresCapability: false,
		ResponseMode:       "text",
		Slots:              map[string]any{},
	}, nil
}

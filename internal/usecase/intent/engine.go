package intent

import (
	"context"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// Engine 是意图识别器的第一版实现
type Engine struct {
}

// NewEngine 创建意图识别器
func NewEngine() *Engine {
	return &Engine{}
}

// Recognize 识别用户意图
func (e *Engine) Recognize(_ context.Context, _ agent.RuntimeContext, message string) (agent.IntentResult, error) {
	lower := strings.ToLower(message)

	// 第一版：规则优先
	switch {
	case strings.Contains(lower, "知识库") || strings.Contains(lower, "检索") || strings.Contains(lower, "rag"):
		return agent.IntentResult{
			IntentType:         agent.IntentRetrievalQA,
			Confidence:         0.90,
			RequiresRAG:        true,
			RequiresPlanning:   false,
			RequiresCapability: false,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, nil

	case strings.Contains(lower, "调用工具") || strings.Contains(lower, "skill") || strings.Contains(lower, "mcp"):
		return agent.IntentResult{
			IntentType:         agent.IntentToolCall,
			Confidence:         0.88,
			RequiresRAG:        false,
			RequiresPlanning:   false,
			RequiresCapability: true,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, nil

	case strings.Contains(lower, "分析") || strings.Contains(lower, "规划") || strings.Contains(lower, "一步一步"):
		return agent.IntentResult{
			IntentType:         agent.IntentWorkflow,
			Confidence:         0.92,
			RequiresRAG:        true,
			RequiresPlanning:   true,
			RequiresCapability: true,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, nil

	case strings.Contains(lower, "写") || strings.Contains(lower, "润色") || strings.Contains(lower, "改写"):
		return agent.IntentResult{
			IntentType:         agent.IntentWrite,
			Confidence:         0.86,
			RequiresRAG:        false,
			RequiresPlanning:   false,
			RequiresCapability: false,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, nil

	default:
		return agent.IntentResult{
			IntentType:         agent.IntentChat,
			Confidence:         0.75,
			RequiresRAG:        false,
			RequiresPlanning:   false,
			RequiresCapability: false,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, nil
	}
}

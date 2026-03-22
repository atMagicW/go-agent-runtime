package intent

import (
	"context"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// RuleClassifier 是基于规则的意图分类器
type RuleClassifier struct {
}

// NewRuleClassifier 创建规则分类器
func NewRuleClassifier() *RuleClassifier {
	return &RuleClassifier{}
}

// Classify 按规则识别意图
func (c *RuleClassifier) Classify(
	_ context.Context,
	_ agent.RuntimeContext,
	message string,
) (agent.IntentResult, bool, error) {
	lower := strings.ToLower(message)

	switch {
	case strings.Contains(lower, "知识库") || strings.Contains(lower, "检索") || strings.Contains(lower, "rag"):
		return agent.IntentResult{
			IntentType:         agent.IntentRetrievalQA,
			Confidence:         0.92,
			RequiresRAG:        true,
			RequiresPlanning:   false,
			RequiresCapability: false,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, true, nil

	case strings.Contains(lower, "调用工具") || strings.Contains(lower, "tool") || strings.Contains(lower, "skill") || strings.Contains(lower, "mcp"):
		return agent.IntentResult{
			IntentType:         agent.IntentToolCall,
			Confidence:         0.90,
			RequiresRAG:        false,
			RequiresPlanning:   false,
			RequiresCapability: true,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, true, nil

	case strings.Contains(lower, "分析") || strings.Contains(lower, "规划") || strings.Contains(lower, "一步一步"):
		return agent.IntentResult{
			IntentType:         agent.IntentWorkflow,
			Confidence:         0.88,
			RequiresRAG:        true,
			RequiresPlanning:   true,
			RequiresCapability: true,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, true, nil

	case strings.Contains(lower, "写") || strings.Contains(lower, "润色") || strings.Contains(lower, "改写"):
		return agent.IntentResult{
			IntentType:         agent.IntentWrite,
			Confidence:         0.86,
			RequiresRAG:        false,
			RequiresPlanning:   false,
			RequiresCapability: false,
			ResponseMode:       "text",
			Slots:              map[string]any{},
		}, true, nil
	}

	return agent.IntentResult{}, false, nil
}

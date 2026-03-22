package intent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// LLMClassifier 是基于模型的意图分类器
type LLMClassifier struct {
	modelRouter ports.ModelRouter
}

// NewLLMClassifier 创建 LLM 分类器
func NewLLMClassifier(modelRouter ports.ModelRouter) *LLMClassifier {
	return &LLMClassifier{
		modelRouter: modelRouter,
	}
}

// llmIntentOutput 表示模型输出结构
type llmIntentOutput struct {
	IntentType         string  `json:"intent_type"`
	Confidence         float64 `json:"confidence"`
	RequiresRAG        bool    `json:"requires_rag"`
	RequiresCapability bool    `json:"requires_capability"`
	RequiresPlanning   bool    `json:"requires_planning"`
	ResponseMode       string  `json:"response_mode"`
}

// Classify 使用模型进行意图识别
func (c *LLMClassifier) Classify(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	message string,
) (agent.IntentResult, bool, error) {
	if c.modelRouter == nil {
		return agent.IntentResult{}, false, nil
	}

	prompt := buildIntentPrompt(message)

	resp, err := c.modelRouter.Generate(ctx, runtimeCtx, ports.ModelCallRequest{
		TaskType: "intent",
		Prompt:   prompt,
		Stream:   false,
	})
	if err != nil {
		return agent.IntentResult{}, false, err
	}

	result, err := parseIntentResult(resp.Text)
	if err != nil {
		return agent.IntentResult{}, false, err
	}

	return result, true, nil
}

func buildIntentPrompt(message string) string {
	return `你是一个 agent 系统的意图分类器。
请根据用户输入识别其意图，并只输出 JSON，不要输出额外解释。

可选意图类型：
- chat
- retrieval_qa
- tool_call
- workflow
- analysis
- write

输出格式：
{
  "intent_type": "chat",
  "confidence": 0.95,
  "requires_rag": false,
  "requires_capability": false,
  "requires_planning": false,
  "response_mode": "text"
}

用户输入：
` + message
}

func parseIntentResult(text string) (agent.IntentResult, error) {
	raw := extractJSONObject(text)
	if raw == "" {
		return agent.IntentResult{}, fmt.Errorf("intent classifier returned invalid json")
	}

	var out llmIntentOutput
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return agent.IntentResult{}, err
	}

	intentType := normalizeIntentType(out.IntentType)

	return agent.IntentResult{
		IntentType:         intentType,
		Confidence:         out.Confidence,
		RequiresRAG:        out.RequiresRAG,
		RequiresCapability: out.RequiresCapability,
		RequiresPlanning:   out.RequiresPlanning,
		ResponseMode:       out.ResponseMode,
		Slots:              map[string]any{},
	}, nil
}

func extractJSONObject(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start < 0 || end < 0 || end <= start {
		return ""
	}
	return text[start : end+1]
}

func normalizeIntentType(s string) agent.IntentType {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "retrieval_qa":
		return agent.IntentRetrievalQA
	case "tool_call":
		return agent.IntentToolCall
	case "workflow":
		return agent.IntentWorkflow
	case "analysis":
		return agent.IntentAnalysis
	case "write":
		return agent.IntentWrite
	default:
		return agent.IntentChat
	}
}

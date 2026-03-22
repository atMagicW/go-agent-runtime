package intent

import (
	"context"
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

type mockModelRouter struct {
	text string
}

func (m *mockModelRouter) Generate(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ModelCallRequest,
) (ports.ModelCallResponse, error) {
	_ = ctx
	_ = runtimeCtx
	_ = req

	return ports.ModelCallResponse{
		Text: m.text,
	}, nil
}

func (m *mockModelRouter) GenerateStream(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ModelCallRequest,
	onToken ports.ModelStreamHandler,
) error {
	_ = ctx
	_ = runtimeCtx
	_ = req
	_ = onToken
	return nil
}

func TestLLMClassifier_Classify(t *testing.T) {
	c := NewLLMClassifier(&mockModelRouter{
		text: `{"intent_type":"workflow","confidence":0.9,"requires_rag":true,"requires_capability":true,"requires_planning":true,"response_mode":"text"}`,
	})

	result, ok, err := c.Classify(context.Background(), agent.RuntimeContext{}, "请分析这个需求并一步一步规划")
	if err != nil {
		t.Fatalf("Classify() error = %v", err)
	}
	if !ok {
		t.Fatal("Classify() ok = false")
	}
	if result.IntentType != agent.IntentWorkflow {
		t.Fatalf("intent type = %s, want workflow", result.IntentType)
	}
}

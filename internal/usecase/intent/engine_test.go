package intent

import (
	"context"
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

func TestRecognize(t *testing.T) {
	e := NewEngine()

	tests := []struct {
		name    string
		message string
		want    agent.IntentType
	}{
		{"rag", "请从知识库检索 agent runtime", agent.IntentRetrievalQA},
		{"tool", "请调用 tool 提取关键词", agent.IntentToolCall},
		{"workflow", "请一步一步分析这个需求", agent.IntentWorkflow},
		{"write", "请帮我润色这段话", agent.IntentWrite},
		{"chat", "你好", agent.IntentChat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := e.Recognize(context.Background(), agent.RuntimeContext{}, tt.message)
			if err != nil {
				t.Fatalf("Recognize() error = %v", err)
			}
			if got.IntentType != tt.want {
				t.Fatalf("Recognize() got = %s, want = %s", got.IntentType, tt.want)
			}
		})
	}
}

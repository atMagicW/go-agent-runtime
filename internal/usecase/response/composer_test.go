package response

import (
	"context"
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	domainprompt "github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

type mockPromptRepo struct{}

func (m *mockPromptRepo) GetByNameAndVersion(ctx context.Context, promptName string, version string) (domainprompt.Template, error) {
	return domainprompt.Template{
		PromptName: promptName,
		Version:    version,
		Content:    "用户请求：{{.Message}}\n意图：{{.Intent}}\n步骤：{{.StepResultsText}}",
	}, nil
}

func (m *mockPromptRepo) GetLatestByName(ctx context.Context, promptName string) (domainprompt.Template, error) {
	return domainprompt.Template{
		PromptName: promptName,
		Version:    "v1",
		Content:    "用户请求：{{.Message}}\n意图：{{.Intent}}\n步骤：{{.StepResultsText}}",
	}, nil
}

func (m *mockPromptRepo) ListByName(ctx context.Context, promptName string) ([]domainprompt.Template, error) {
	return []domainprompt.Template{
		{
			PromptName: promptName,
			Version:    "v1",
			Content:    "用户请求：{{.Message}}",
		},
	}, nil
}

var _ ports.PromptRepository = (*mockPromptRepo)(nil)

func TestCompose(t *testing.T) {
	c := NewTemplateResponseComposer(&mockPromptRepo{})

	resp, err := c.Compose(context.Background(), agent.RuntimeContext{
		Intent: agent.IntentResult{
			IntentType: agent.IntentWorkflow,
		},
	}, ports.ComposeRequest{
		Message:    "请分析这个需求",
		PromptName: "final_response",
		StepResults: []agent.StepResult{
			{
				StepID:  "step_1",
				Success: true,
				Output: map[string]any{
					"result": "步骤执行成功",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Compose error = %v", err)
	}

	if resp.Text == "" {
		t.Fatal("response text is empty")
	}
	t.Log(resp.Text)
}

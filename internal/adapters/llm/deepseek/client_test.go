package deepseek

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

type testPricing struct{}

func (testPricing) CalcLLMCost(model string, inputTokens, outputTokens int) float64 {
	return 0.0
}

func getTestClient(t *testing.T) *Client {
	t.Helper()

	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey == "" {
		t.Skip("skip integration test: DEEPSEEK_API_KEY is not set")
	}

	baseURL := strings.TrimSpace(os.Getenv("DEEPSEEK_BASE_URL"))
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}

	return NewClient(apiKey, baseURL, testPricing{})
}

func TestClient_Generate(t *testing.T) {
	client := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, ports.LLMGenerateRequest{
		Prompt: "请只回复：hello",
		Model:  "deepseek-chat",
	})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if strings.TrimSpace(resp.Text) == "" {
		t.Fatal("Generate returned empty text")
	}

	if resp.Model != "deepseek-chat" {
		t.Fatalf("unexpected model: %s", resp.Model)
	}

	if resp.Provider != "deepseek" {
		t.Fatalf("unexpected provider: %s", resp.Provider)
	}

	if resp.PromptTokens < 0 || resp.CompletionTokens < 0 || resp.TotalTokens < 0 {
		t.Fatalf("invalid token stats: %+v", resp)
	}

	if resp.TotalTokens == 0 {
		t.Fatalf("TotalTokens should not be 0, got %+v", resp)
	}

	t.Logf("Generate text: %q", resp.Text)
	t.Logf("usage: prompt=%d completion=%d total=%d",
		resp.PromptTokens, resp.CompletionTokens, resp.TotalTokens)
}

func TestClient_Generate_EmptyPrompt(t *testing.T) {
	client := &Client{}

	_, err := client.Generate(context.Background(), ports.LLMGenerateRequest{
		Prompt: "",
		Model:  "deepseek-chat",
	})
	if err == nil {
		t.Fatal("expected error for empty prompt, got nil")
	}
}

func TestClient_GenerateStream(t *testing.T) {
	client := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var chunks []string
	doneCalled := false

	err := client.GenerateStream(ctx, ports.LLMGenerateRequest{
		Prompt: "请用一句很短的话介绍 agent。",
		Model:  "deepseek-chat",
	}, func(chunk ports.StreamChunk) error {
		if chunk.Done {
			doneCalled = true
			return nil
		}
		if chunk.Text != "" {
			chunks = append(chunks, chunk.Text)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("GenerateStream failed: %v", err)
	}

	fullText := strings.TrimSpace(strings.Join(chunks, ""))
	if fullText == "" {
		t.Fatal("GenerateStream returned empty content")
	}

	if !doneCalled {
		t.Fatal("GenerateStream did not send Done=true chunk")
	}

	t.Logf("GenerateStream text: %q", fullText)
}

func TestClient_GenerateStream_EmptyPrompt(t *testing.T) {
	client := &Client{}

	err := client.GenerateStream(context.Background(), ports.LLMGenerateRequest{
		Prompt: "",
		Model:  "deepseek-chat",
	}, func(chunk ports.StreamChunk) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error for empty prompt, got nil")
	}
}

func TestClient_GenerateStream_NilHandler(t *testing.T) {
	client := &Client{}

	err := client.GenerateStream(context.Background(), ports.LLMGenerateRequest{
		Prompt: "hello",
		Model:  "deepseek-chat",
	}, nil)
	if err == nil {
		t.Fatal("expected error for nil stream handler, got nil")
	}
}

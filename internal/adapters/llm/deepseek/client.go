package deepseek

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

type CostCalculator interface {
	CalcLLMCost(model string, inputTokens, outputTokens int) float64
}

// Client 是 DeepSeek 的适配器实现。
// DeepSeek API 与 OpenAI 兼容，因此这里继续复用 OpenAI Go SDK。
type Client struct {
	client  openai.Client
	pricing CostCalculator
	baseURL string
}

// NewClient 创建 DeepSeek 客户端
func NewClient(apiKey string, baseURL string, pricing CostCalculator) *Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.deepseek.com"
	}

	c := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	return &Client{
		client:  c,
		pricing: pricing,
		baseURL: baseURL,
	}
}

// Generate 非流式生成
func (c *Client) Generate(ctx context.Context, req ports.LLMGenerateRequest) (ports.LLMGenerateResponse, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return ports.LLMGenerateResponse{}, fmt.Errorf("prompt is empty")
	}

	modelName := req.Model
	if strings.TrimSpace(modelName) == "" {
		modelName = "deepseek-chat"
	}

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: modelName,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(req.Prompt),
		},
	})
	if err != nil {
		return ports.LLMGenerateResponse{}, err
	}

	text := ""
	if len(resp.Choices) > 0 {
		text = resp.Choices[0].Message.Content
	}

	promptTokens := 0
	completionTokens := 0
	totalTokens := 0

	if resp.Usage.PromptTokens > 0 {
		promptTokens = int(resp.Usage.PromptTokens)
	}
	if resp.Usage.CompletionTokens > 0 {
		completionTokens = int(resp.Usage.CompletionTokens)
	}
	if resp.Usage.TotalTokens > 0 {
		totalTokens = int(resp.Usage.TotalTokens)
	}

	if promptTokens == 0 {
		promptTokens = estimateTokens(req.Prompt)
	}
	if completionTokens == 0 {
		completionTokens = estimateTokens(text)
	}
	if totalTokens == 0 {
		totalTokens = promptTokens + completionTokens
	}

	cost := 0.0
	if c.pricing != nil {
		cost = c.pricing.CalcLLMCost(modelName, promptTokens, completionTokens)
	}

	return ports.LLMGenerateResponse{
		Text:             text,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokens,
		Cost:             cost,
		Model:            modelName,
		Provider:         "deepseek",
	}, nil
}

// GenerateStream 流式生成
func (c *Client) GenerateStream(ctx context.Context, req ports.LLMGenerateRequest, onChunk ports.StreamHandler) error {
	if strings.TrimSpace(req.Prompt) == "" {
		return fmt.Errorf("prompt is empty")
	}
	if onChunk == nil {
		return fmt.Errorf("stream handler is nil")
	}

	modelName := req.Model
	if strings.TrimSpace(modelName) == "" {
		modelName = "deepseek-chat"
	}

	stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(modelName),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(req.Prompt),
		},
	})

	defer stream.Close()

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta.Content
		if delta == "" {
			continue
		}

		if err := onChunk(ports.StreamChunk{Text: delta}); err != nil {
			return err
		}
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return onChunk(ports.StreamChunk{Done: true})
}

func estimateTokens(text string) int {
	if text == "" {
		return 0
	}
	n := len([]rune(text)) / 4
	if n == 0 {
		return 1
	}
	return n
}

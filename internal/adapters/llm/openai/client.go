package openai

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"

	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

type CostCalculator interface {
	CalcLLMCost(model string, inputTokens, outputTokens int) float64
}

type Client struct {
	client  openai.Client
	pricing CostCalculator
}

func NewClient(apiKey string, pricing CostCalculator) *Client {
	var c openai.Client
	if strings.TrimSpace(apiKey) != "" {
		c = openai.NewClient(option.WithAPIKey(apiKey))
	} else {
		c = openai.NewClient()
	}

	return &Client{
		client:  c,
		pricing: pricing,
	}
}

func (c *Client) Generate(ctx context.Context, req ports.LLMGenerateRequest) (ports.LLMGenerateResponse, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return ports.LLMGenerateResponse{}, fmt.Errorf("prompt is empty")
	}

	modelName := req.Model
	if strings.TrimSpace(modelName) == "" {
		modelName = "gpt-4.1-mini"
	}

	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(req.Prompt),
		},
		Model: openai.ChatModel(modelName),
	})
	if err != nil {
		return ports.LLMGenerateResponse{}, err
	}

	text := resp.OutputText()

	// 优先使用 SDK 返回的 usage；如果本地版本字段不同，只改这一处即可。
	inputTokens := 0
	outputTokens := 0
	totalTokens := 0

	// 当前官方 SDK 文档/类型中有 ResponseUsage。实际字段名若随版本有细微变化，保留这一层适配。:contentReference[oaicite:2]{index=2}
	if resp.Usage.InputTokens > 0 {
		inputTokens = int(resp.Usage.InputTokens)
	}
	if resp.Usage.OutputTokens > 0 {
		outputTokens = int(resp.Usage.OutputTokens)
	}
	if resp.Usage.TotalTokens > 0 {
		totalTokens = int(resp.Usage.TotalTokens)
	}

	// SDK usage 取不到时再退回估算
	if inputTokens == 0 {
		inputTokens = estimateTokens(req.Prompt)
	}
	if outputTokens == 0 {
		outputTokens = estimateTokens(text)
	}
	if totalTokens == 0 {
		totalTokens = inputTokens + outputTokens
	}

	cost := 0.0
	if c.pricing != nil {
		cost = c.pricing.CalcLLMCost(modelName, inputTokens, outputTokens)
	}

	return ports.LLMGenerateResponse{
		Text:             text,
		PromptTokens:     inputTokens,
		CompletionTokens: outputTokens,
		TotalTokens:      totalTokens,
		Cost:             cost,
		Model:            modelName,
		Provider:         "openai",
	}, nil
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

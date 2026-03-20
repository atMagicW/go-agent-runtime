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

// Client 是 OpenAI 的适配器实现
type Client struct {
	client openai.Client
}

// NewClient 创建 OpenAI 客户端
func NewClient(apiKey string) *Client {
	var c openai.Client
	if strings.TrimSpace(apiKey) != "" {
		c = openai.NewClient(option.WithAPIKey(apiKey))
	} else {
		c = openai.NewClient()
	}

	return &Client{
		client: c,
	}
}

// Generate 同步生成
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
	promptTokens := estimateTokens(req.Prompt)
	completionTokens := estimateTokens(text)

	return ports.LLMGenerateResponse{
		Text:             text,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		Cost:             estimateOpenAICost(modelName, promptTokens, completionTokens),
		Model:            modelName,
		Provider:         "openai",
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
		modelName = "gpt-4.1-mini"
	}

	stream := c.client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(req.Prompt),
		},
		Model: openai.ChatModel(modelName),
	})
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()

		// 第一版先宽松处理：
		// 只要事件里能提取到增量文本，就推给上层。
		// 所以这里用字符串化兜底不太好，优先走常见文本增量字段。
		if delta := extractTextDelta(event); delta != "" {
			if err := onChunk(ports.StreamChunk{Text: delta}); err != nil {
				return err
			}
		}
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return onChunk(ports.StreamChunk{Done: true})
}

// extractTextDelta 从 streaming event 中提取文本增量
func extractTextDelta(event responses.ResponseStreamEventUnion) string {
	// 这里根据 SDK 常见 Responses streaming 事件做兼容。
	// 如果你本地 SDK 版本字段略有差异，只需要改这里。
	if event.Type == "response.output_text.delta" {
		return event.AsResponseOutputTextDelta().Delta
	}
	return ""
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

func estimateOpenAICost(model string, promptTokens, completionTokens int) float64 {
	switch model {
	case "gpt-4.1":
		return float64(promptTokens+completionTokens) * 0.00001
	case "gpt-4.1-mini":
		return float64(promptTokens+completionTokens) * 0.000002
	default:
		return float64(promptTokens+completionTokens) * 0.000003
	}
}

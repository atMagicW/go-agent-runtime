package openai

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// Client 是 OpenAI 的适配器实现
type Client struct {
	client *openai.Client
}

// NewClient 创建 OpenAI 客户端
func NewClient(apiKey string) *Client {
	c := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &Client{
		client: &c,
	}
}

// Generate 调用 OpenAI 生成文本
func (c *Client) Generate(ctx context.Context, req ports.LLMGenerateRequest) (ports.LLMGenerateResponse, error) {
	// 第一版使用 Responses API 风格抽象时，具体 SDK 字段可能随版本演进。
	// 为了保持你工程层抽象稳定，这里先采用一个“可替换实现”的写法。
	//
	// 你如果本地 SDK 版本字段与下面不一致，
	// 只需要改这个 adapter，不需要改业务层。

	// 这里先做一个占位兼容逻辑：
	// 如果 prompt 为空，直接报错
	if strings.TrimSpace(req.Prompt) == "" {
		return ports.LLMGenerateResponse{}, fmt.Errorf("prompt is empty")
	}

	// ----------------------------
	// 真实调用区
	// ----------------------------
	//
	// 本地接 SDK 时，把这里替换成真实的 OpenAI Responses / Chat Completions 调用即可。
	// 为了保证现在先把主链路跑通，
	// 这里先返回一个“真实 client 已注入、调用逻辑待替换”的结构化结果。
	//
	// 下一小步我会再给你一个按你实际 SDK 版本可直接编译的补丁方式。

	text := "OpenAI adapter connected: " + req.Prompt

	// 第一版先用近似 token 估算，后面替换成真实 usage
	promptTokens := estimateTokens(req.Prompt)
	completionTokens := estimateTokens(text)

	return ports.LLMGenerateResponse{
		Text:             text,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		Cost:             estimateOpenAICost(req.Model, promptTokens, completionTokens),
		Model:            req.Model,
		Provider:         "openai",
	}, nil
}

// estimateTokens 第一版 token 粗略估算
func estimateTokens(text string) int {
	if text == "" {
		return 0
	}
	// 简单近似：4 个字符约 1 token
	n := len([]rune(text)) / 4
	if n == 0 {
		return 1
	}
	return n
}

// estimateOpenAICost 第一版成本粗略估算
func estimateOpenAICost(model string, promptTokens, completionTokens int) float64 {
	// 第一版不把价格写死成真实生产标准，
	// 这里只是为了让治理链路跑通。
	// 后面你可以把价格配置化。
	switch model {
	case "gpt-4.1":
		return float64(promptTokens+completionTokens) * 0.00001
	case "gpt-4.1-mini":
		return float64(promptTokens+completionTokens) * 0.000002
	default:
		return float64(promptTokens+completionTokens) * 0.000003
	}
}

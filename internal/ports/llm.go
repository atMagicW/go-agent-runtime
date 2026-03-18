package ports

import "context"

// LLMGenerateRequest 表示一次标准化模型生成请求
type LLMGenerateRequest struct {
	Model  string
	Prompt string
	Stream bool
}

// LLMGenerateResponse 表示一次标准化模型生成结果
type LLMGenerateResponse struct {
	Text string

	PromptTokens     int
	CompletionTokens int
	TotalTokens      int

	Cost float64

	Model    string
	Provider string
}

// LLMClient 定义统一的模型客户端接口
type LLMClient interface {
	Generate(ctx context.Context, req LLMGenerateRequest) (LLMGenerateResponse, error)
}

package ports

import "context"

// StreamChunk 表示一次流式输出片段
type StreamChunk struct {
	Text string
	Done bool
}

// StreamHandler 处理流式输出片段
type StreamHandler func(chunk StreamChunk) error

// StreamingLLMClient 定义支持流式输出的模型客户端接口
type StreamingLLMClient interface {
	GenerateStream(ctx context.Context, req LLMGenerateRequest, onChunk StreamHandler) error
}

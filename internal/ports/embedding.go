package ports

import "context"

// EmbeddingProvider 定义向量化接口
type EmbeddingProvider interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

package mockembedding

import (
	"context"
	"hash/fnv"
)

// Provider 是一个可重复的 mock embedding provider
type Provider struct {
	dim int
}

// NewProvider 创建 mock embedding provider
func NewProvider(dim int) *Provider {
	if dim <= 0 {
		dim = 1536
	}
	return &Provider{dim: dim}
}

// Embed 生成稳定的伪向量
func (p *Provider) Embed(_ context.Context, text string) ([]float32, error) {
	vec := make([]float32, p.dim)

	h := fnv.New32a()
	_, _ = h.Write([]byte(text))
	base := h.Sum32()

	for i := 0; i < p.dim; i++ {
		v := float32((int(base)+i*31)%1000) / 1000.0
		vec[i] = v
	}

	return vec, nil
}

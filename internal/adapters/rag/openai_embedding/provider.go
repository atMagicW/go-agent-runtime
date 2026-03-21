package openaiembedding

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type CostCalculator interface {
	CalcEmbeddingCost(model string, inputTokens int) float64
}

type Provider struct {
	client  openai.Client
	model   string
	dim     int
	pricing CostCalculator
}

func NewProvider(apiKey string, model string, dim int, pricing CostCalculator) *Provider {
	var c openai.Client
	if strings.TrimSpace(apiKey) != "" {
		c = openai.NewClient(option.WithAPIKey(apiKey))
	} else {
		c = openai.NewClient()
	}

	if model == "" {
		model = "text-embedding-3-small"
	}

	return &Provider{
		client:  c,
		model:   model,
		dim:     dim,
		pricing: pricing,
	}
}

// Embed 将文本转换为向量
func (p *Provider) Embed(ctx context.Context, text string) ([]float32, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("text is empty")
	}

	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
		Model: openai.EmbeddingModel(p.model),
	}

	// 仅 text-embedding-3 系列支持 dimensions 参数；第一版先按配置尝试设置。
	if p.dim > 0 && strings.HasPrefix(p.model, "text-embedding-3") {
		params.Dimensions = openai.Int(int64(p.dim))
	}

	resp, err := p.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("embedding response is empty")
	}

	raw := resp.Data[0].Embedding
	out := make([]float32, len(raw))
	for i, v := range raw {
		out[i] = float32(v)
	}

	// 这里如果你后面要把 embedding cost 也落库，
	// 可以把 resp.Usage.PromptTokens 往外透出。
	_ = resp

	return out, nil
}

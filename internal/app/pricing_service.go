package app

import cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"

type PricingService struct {
	cfg *cfg.PricingConfig
}

func NewPricingService(c *cfg.PricingConfig) *PricingService {
	return &PricingService{cfg: c}
}

func (s *PricingService) CalcLLMCost(model string, inputTokens, outputTokens int) float64 {
	if s == nil || s.cfg == nil {
		return 0
	}

	price, ok := s.cfg.LLMPricing[model]
	if !ok {
		return 0
	}

	inputCost := (float64(inputTokens) / 1_000_000.0) * price.InputPerMillion
	outputCost := (float64(outputTokens) / 1_000_000.0) * price.OutputPerMillion
	return inputCost + outputCost
}

func (s *PricingService) CalcEmbeddingCost(model string, inputTokens int) float64 {
	if s == nil || s.cfg == nil {
		return 0
	}

	price, ok := s.cfg.EmbeddingPricing[model]
	if !ok {
		return 0
	}

	return (float64(inputTokens) / 1_000_000.0) * price.InputPerMillion
}

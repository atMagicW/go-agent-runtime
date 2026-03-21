package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type LLMPrice struct {
	InputPerMillion  float64 `yaml:"input_per_million"`
	OutputPerMillion float64 `yaml:"output_per_million"`
}

type EmbeddingPrice struct {
	InputPerMillion float64 `yaml:"input_per_million"`
}

type PricingConfig struct {
	LLMPricing       map[string]LLMPrice       `yaml:"llm_pricing"`
	EmbeddingPricing map[string]EmbeddingPrice `yaml:"embedding_pricing"`
}

func LoadPricing(path string) (*PricingConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pricing config failed: %w", err)
	}

	var cfg PricingConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal pricing config failed: %w", err)
	}

	if cfg.LLMPricing == nil {
		cfg.LLMPricing = map[string]LLMPrice{}
	}
	if cfg.EmbeddingPricing == nil {
		cfg.EmbeddingPricing = map[string]EmbeddingPrice{}
	}

	return &cfg, nil
}

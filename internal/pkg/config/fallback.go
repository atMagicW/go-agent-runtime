package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FallbackConfig 表示 fallback 配置文件
type FallbackConfig struct {
	ModelFallbacks         map[string][]string `yaml:"model_fallbacks"`
	CapabilityFallbacks    map[string][]string `yaml:"capability_fallbacks"`
	KnowledgeBaseFallbacks map[string][]string `yaml:"knowledge_base_fallbacks"`
}

// LoadFallback 加载 fallback 配置
func LoadFallback(path string) (*FallbackConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read fallback config failed: %w", err)
	}

	var cfg FallbackConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal fallback config failed: %w", err)
	}

	if cfg.ModelFallbacks == nil {
		cfg.ModelFallbacks = map[string][]string{}
	}
	if cfg.CapabilityFallbacks == nil {
		cfg.CapabilityFallbacks = map[string][]string{}
	}
	if cfg.KnowledgeBaseFallbacks == nil {
		cfg.KnowledgeBaseFallbacks = map[string][]string{}
	}

	return &cfg, nil
}

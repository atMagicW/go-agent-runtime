package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ModelConfig 表示单个模型配置
type ModelConfig struct {
	Name     string   `yaml:"name"`
	Provider string   `yaml:"provider"`
	Enabled  bool     `yaml:"enabled"`
	Tags     []string `yaml:"tags"`
}

// ModelsConfig 表示模型配置文件
type ModelsConfig struct {
	DefaultModel string        `yaml:"default_model"`
	Models       []ModelConfig `yaml:"models"`
}

// LoadModels 加载模型配置
func LoadModels(path string) (*ModelsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read models config failed: %w", err)
	}

	var cfg ModelsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal models config failed: %w", err)
	}

	return &cfg, nil
}

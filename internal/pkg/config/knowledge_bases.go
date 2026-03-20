package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// KnowledgeBaseConfig 表示知识库配置
type KnowledgeBaseConfig struct {
	KBID         string `yaml:"kb_id"`
	TenantID     string `yaml:"tenant_id"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Enabled      bool   `yaml:"enabled"`
	SeedDemoData bool   `yaml:"seed_demo_data"`
}

// KnowledgeBasesConfig 表示知识库配置文件
type KnowledgeBasesConfig struct {
	KnowledgeBases []KnowledgeBaseConfig `yaml:"knowledge_bases"`
}

// LoadKnowledgeBases 加载知识库配置
func LoadKnowledgeBases(path string) (*KnowledgeBasesConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read knowledge bases config failed: %w", err)
	}

	var cfg KnowledgeBasesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal knowledge bases config failed: %w", err)
	}

	return &cfg, nil
}

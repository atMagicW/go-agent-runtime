package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 表示应用主配置
type Config struct {
	App struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
		Port int    `yaml:"port"`
	} `yaml:"app"`

	Database struct {
		PostgresDSN string `yaml:"postgres_dsn"`
	} `yaml:"database"`

	LLM struct {
		OpenAIAPIKey string `yaml:"openai_api_key"`
		DefaultModel string `yaml:"default_model"`
	} `yaml:"llm"`

	RAG struct {
		DefaultTopK       int    `yaml:"default_top_k"`
		EmbeddingDim      int    `yaml:"embedding_dim"`
		RerankEnabled     bool   `yaml:"rerank_enabled"`
		EmbeddingProvider string `yaml:"embedding_provider"`
		EmbeddingModel    string `yaml:"embedding_model"`
	} `yaml:"rag"`

	TextSplitter struct {
		ChunkSize int `yaml:"chunk_size"`
		Overlap   int `yaml:"overlap"`
	} `yaml:"text_splitter"`
}

// Load 从 yaml 文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	applyEnvOverrides(&cfg)
	applyDefaults(&cfg)

	return &cfg, nil
}

// applyEnvOverrides 使用环境变量覆盖配置
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("POSTGRES_DSN"); v != "" {
		cfg.Database.PostgresDSN = v
	}

	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.LLM.OpenAIAPIKey = v
	}

	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.App.Env = v
	}
}

// applyDefaults 填默认值
func applyDefaults(cfg *Config) {
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}

	if cfg.LLM.DefaultModel == "" {
		cfg.LLM.DefaultModel = "gpt-4.1-mini"
	}

	if cfg.RAG.DefaultTopK <= 0 {
		cfg.RAG.DefaultTopK = 5
	}

	if cfg.RAG.EmbeddingDim <= 0 {
		cfg.RAG.EmbeddingDim = 1536
	}

	if cfg.RAG.EmbeddingProvider == "" {
		cfg.RAG.EmbeddingProvider = "openai"
	}

	if cfg.RAG.EmbeddingModel == "" {
		cfg.RAG.EmbeddingModel = "text-embedding-3-small"
	}

	if cfg.TextSplitter.ChunkSize <= 0 {
		cfg.TextSplitter.ChunkSize = 300
	}

	if cfg.TextSplitter.Overlap < 0 {
		cfg.TextSplitter.Overlap = 50
	}
}

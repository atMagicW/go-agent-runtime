package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SimpleCapabilityConfig 表示本地 skill/tool 的开关配置
type SimpleCapabilityConfig struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
}

// MCPToolConfig 表示 MCP Tool 配置
type MCPToolConfig struct {
	Name        string `yaml:"name"`
	ServerName  string `yaml:"server_name"`
	RemoteTool  string `yaml:"remote_tool"`
	Description string `yaml:"description"`
	Enabled     bool   `yaml:"enabled"`
}

// CapabilitiesConfig 表示能力配置文件
type CapabilitiesConfig struct {
	Skills   []SimpleCapabilityConfig `yaml:"skills"`
	Tools    []SimpleCapabilityConfig `yaml:"tools"`
	MCPTools []MCPToolConfig          `yaml:"mcp_tools"`
}

// LoadCapabilities 加载能力配置
func LoadCapabilities(path string) (*CapabilitiesConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read capabilities config failed: %w", err)
	}

	var cfg CapabilitiesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal capabilities config failed: %w", err)
	}

	return &cfg, nil
}

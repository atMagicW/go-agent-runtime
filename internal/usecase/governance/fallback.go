package governance

import cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"

// FallbackPolicy 定义简单降级策略
type FallbackPolicy struct {
	ModelFallbacks         map[string][]string
	CapabilityFallbacks    map[string][]string
	KnowledgeBaseFallbacks map[string][]string
}

// NewDefaultFallbackPolicy 创建默认降级策略
func NewDefaultFallbackPolicy() *FallbackPolicy {
	return &FallbackPolicy{
		ModelFallbacks: map[string][]string{
			"gpt-4.1":      {"gpt-4.1-mini"},
			"gpt-4.1-mini": {"gpt-4.1-mini"},
		},
		CapabilityFallbacks: map[string][]string{
			"mcp_web_search": {"keyword_extract_tool"},
		},
		KnowledgeBaseFallbacks: map[string][]string{
			"knowledge_a": {"default"},
			"knowledge_b": {"default"},
		},
	}
}

// NewFallbackPolicyFromConfig 根据配置构造降级策略
func NewFallbackPolicyFromConfig(c *cfg.FallbackConfig) *FallbackPolicy {
	if c == nil {
		return NewDefaultFallbackPolicy()
	}

	return &FallbackPolicy{
		ModelFallbacks:         c.ModelFallbacks,
		CapabilityFallbacks:    c.CapabilityFallbacks,
		KnowledgeBaseFallbacks: c.KnowledgeBaseFallbacks,
	}
}

func (p *FallbackPolicy) NextModels(model string) []string {
	if p == nil {
		return nil
	}
	return p.ModelFallbacks[model]
}

// NextCapabilities 获取能力回退链
func (p *FallbackPolicy) NextCapabilities(name string) []string {
	if p == nil {
		return nil
	}
	return p.CapabilityFallbacks[name]
}

// NextKnowledgeBases 获取知识库回退链
func (p *FallbackPolicy) NextKnowledgeBases(name string) []string {
	if p == nil {
		return nil
	}
	return p.KnowledgeBaseFallbacks[name]
}

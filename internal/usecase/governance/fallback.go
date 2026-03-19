package governance

// FallbackPolicy 定义简单降级策略
type FallbackPolicy struct {
	// 模型回退链
	ModelFallbacks map[string][]string

	// capability 回退链
	CapabilityFallbacks map[string][]string

	// kb 回退链
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

// NextModels 获取模型回退链
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

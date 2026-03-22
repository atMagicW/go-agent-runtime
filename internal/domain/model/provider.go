package model

// ProviderName 表示模型提供商名称
type ProviderName string

const (
	ProviderOpenAI    ProviderName = "openai"
	ProviderDeepSeek  ProviderName = "deepseek"
	ProviderAnthropic ProviderName = "anthropic"
	ProviderGemini    ProviderName = "gemini"
)

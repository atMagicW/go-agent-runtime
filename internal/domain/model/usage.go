package model

import "time"

// UsageRecord 表示一次模型调用的用量记录
type UsageRecord struct {
	RequestID string `json:"request_id"`
	SessionID string `json:"session_id"`

	Provider string `json:"provider"`
	Model    string `json:"model"`

	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`

	// 本次调用成本，单位先统一按美元计算
	Cost float64 `json:"cost"`

	LatencyMS int64 `json:"latency_ms"`

	CreatedAt time.Time `json:"created_at"`
}

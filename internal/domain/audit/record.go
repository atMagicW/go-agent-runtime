package audit

import "time"

// Record 表示一次请求的审计记录
type Record struct {
	AuditID string `json:"audit_id"`

	RequestID string `json:"request_id"`
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`

	Intent string `json:"intent"`

	// 执行计划快照，第一版直接存字符串
	PlanSnapshot string `json:"plan_snapshot"`

	// 使用到的模型
	ModelsUsed []string `json:"models_used,omitempty"`

	// 使用到的知识库
	KnowledgeBases []string `json:"knowledge_bases,omitempty"`

	// 调用到的能力
	Capabilities []string `json:"capabilities,omitempty"`

	// 最终结果摘要
	ResultSummary string `json:"result_summary"`

	Status string `json:"status"`

	DurationMS int64 `json:"duration_ms"`

	TotalCost float64 `json:"total_cost"`

	CreatedAt time.Time `json:"created_at"`
}

package agent

import "time"

// StepResult 表示单个步骤的执行结果
type StepResult struct {
	StepID string `json:"step_id"`

	Success bool `json:"success"`

	// 标准化输出
	Output map[string]any `json:"output,omitempty"`

	// 错误信息
	Error string `json:"error,omitempty"`

	StartedAt time.Time `json:"started_at"`

	EndedAt time.Time `json:"ended_at"`
}

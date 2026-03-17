package agent

import "time"

// PlanStatus 表示执行计划状态
type PlanStatus string

const (
	PlanStatusPending   PlanStatus = "pending"
	PlanStatusRunning   PlanStatus = "running"
	PlanStatusSucceeded PlanStatus = "succeeded"
	PlanStatusFailed    PlanStatus = "failed"
)

// ExecutionPlan 表示一次任务的执行计划
type ExecutionPlan struct {
	PlanID string `json:"plan_id"`

	// 用户目标
	Goal string `json:"goal"`

	// 计划状态
	Status PlanStatus `json:"status"`

	// 步骤列表
	Steps []PlanStep `json:"steps"`

	CreatedAt time.Time `json:"created_at"`
}

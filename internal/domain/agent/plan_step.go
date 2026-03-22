package agent

// StepType 表示步骤类型
type StepType string

const (
	StepTypeIntentRefine    StepType = "intent_refine"
	StepTypeRetrieve        StepType = "retrieve"
	StepTypeTool            StepType = "tool"
	StepTypeMCP             StepType = "mcp"
	StepTypeLLMGenerate     StepType = "llm_generate"
	StepTypeLLMAnalyze      StepType = "llm_analyze"
	StepTypeSummarize       StepType = "summarize"
	StepTypeComposeResponse StepType = "compose_response"
)

// StepStatus 表示步骤执行状态
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusSucceeded StepStatus = "succeeded"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
)

// RetryPolicy 表示步骤重试策略
type RetryPolicy struct {
	MaxRetries int `json:"max_retries"`
}

// PlanStep 表示执行计划中的一个步骤
type PlanStep struct {
	StepID string `json:"step_id"`

	Name string `json:"name"`

	Type StepType `json:"type"`

	// 依赖的前置步骤 ID
	DependsOn []string `json:"depends_on,omitempty"`

	// 并行分组，相同 group 且无依赖冲突时可并发
	ParallelGroup string `json:"parallel_group,omitempty"`

	// 执行器名称，例如 model_router / rag_router / capability_router
	Executor ExecutorName `json:"executor"`

	// 输入参数
	Input map[string]any `json:"input,omitempty"`

	// 步骤级超时，单位毫秒
	TimeoutMS int64 `json:"timeout_ms,omitempty"`

	// 重试策略
	RetryPolicy RetryPolicy `json:"retry_policy,omitempty"`

	// 失败后的回退动作，例如 fallback_model / fallback_keyword_retrieval
	Fallback string `json:"fallback,omitempty"`

	Status StepStatus `json:"status"`
}

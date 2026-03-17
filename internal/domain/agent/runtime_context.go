package agent

// RuntimeContext 表示本次 Agent 执行使用的运行时上下文
type RuntimeContext struct {
	Request RequestContext `json:"request"`

	Conversation ConversationState `json:"conversation"`

	Intent IntentResult `json:"intent"`

	// 当前已产生的中间变量
	Variables map[string]any `json:"variables,omitempty"`
}

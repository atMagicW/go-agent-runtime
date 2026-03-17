package agent

// IntentType 表示识别出的意图类型
type IntentType string

const (
	IntentChat            IntentType = "chat"
	IntentRetrievalQA     IntentType = "retrieval_qa"
	IntentToolCall        IntentType = "tool_call"
	IntentWorkflow        IntentType = "workflow"
	IntentAnalysis        IntentType = "analysis"
	IntentWrite           IntentType = "write"
	IntentAgenticResearch IntentType = "agentic_research"
	IntentOperate         IntentType = "operate"
)

// IntentResult 表示意图识别结果
type IntentResult struct {
	// 意图类型
	IntentType IntentType `json:"intent_type"`

	// 置信度
	Confidence float64 `json:"confidence"`

	// 槽位信息，第一版先用 map 承载
	Slots map[string]any `json:"slots,omitempty"`

	// 是否需要检索
	RequiresRAG bool `json:"requires_rag"`

	// 是否需要调用 Skill / MCP
	RequiresCapability bool `json:"requires_capability"`

	// 是否需要进入任务规划
	RequiresPlanning bool `json:"requires_planning"`

	// 回复模式，例如 text / stream / json
	ResponseMode string `json:"response_mode,omitempty"`
}

package prompt

// PromptName 表示系统内置 prompt 名称
type PromptName string

const (
	PromptFinalResponse    PromptName = "final_response"
	PromptIntentClassifier PromptName = "intent_classifier"
	PromptMakePlan         PromptName = "make_plan"
)

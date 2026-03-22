package agent

// ExecutorName 表示步骤执行器名称
type ExecutorName string

const (
	ExecutorModelRouter      ExecutorName = "model_router"
	ExecutorCapabilityRouter ExecutorName = "capability_router"
	ExecutorRAGRouter        ExecutorName = "rag_router"
	ExecutorResponseComposer ExecutorName = "response_composer"
)

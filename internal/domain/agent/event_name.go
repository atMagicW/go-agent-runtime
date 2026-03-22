package agent

// EventName 表示运行时流式事件名称
type EventName string

const (
	EventIntent              EventName = "intent"
	EventPlan                EventName = "plan"
	EventStepStarted         EventName = "step_started"
	EventStepCompleted       EventName = "step_completed"
	EventStepFailed          EventName = "step_failed"
	EventRetrieval           EventName = "retrieval"
	EventToolCalled          EventName = "tool_called"
	EventFinalAnswer         EventName = "final_answer"
	EventFinalAnswerStart    EventName = "final_answer_start"
	EventFinalAnswerFallback EventName = "final_answer_fallback"
	EventDone                EventName = "done"
	EventError               EventName = "error"
	EventToken               EventName = "token"
)

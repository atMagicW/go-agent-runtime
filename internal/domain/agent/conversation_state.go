package agent

import "time"

// Message 表示一条会话消息
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ConversationState 表示多轮会话状态
type ConversationState struct {
	SessionID string `json:"session_id"`

	Messages []Message `json:"messages,omitempty"`

	// 历史摘要
	Summary string `json:"summary,omitempty"`

	// 当前活跃任务 ID
	ActiveTaskID string `json:"active_task_id,omitempty"`

	// 变量快照
	Variables map[string]any `json:"variables,omitempty"`
}

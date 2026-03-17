package agent

import "time"

// Session 表示一个持久化会话
type Session struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`

	Summary string `json:"summary,omitempty"`

	ActiveTaskID string `json:"active_task_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

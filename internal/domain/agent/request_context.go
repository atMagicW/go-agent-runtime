package agent

import "time"

// RequestContext 表示一次Agent请求的上下文
type RequestContext struct {

	// 请求ID
	RequestID string

	// 会话ID
	SessionID string

	// 用户ID
	UserID string

	// 用户指定模型
	Model string

	// 是否流式输出
	Stream bool

	// 请求超时时间
	Deadline time.Time
}

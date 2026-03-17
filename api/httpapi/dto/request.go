package httpapi

// ChatRequest 表示聊天请求
type ChatRequest struct {
	// 会话 ID，前端可传入，服务端也可以自动生成
	SessionID string `json:"session_id"`

	// 用户 ID，用于多用户隔离
	UserID string `json:"user_id"`

	// 用户输入内容
	Message string `json:"message"`

	// 用户指定模型，可为空
	Model string `json:"model"`

	// 是否启用流式输出
	Stream bool `json:"stream"`
}

package httpapi

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// StreamWriter 用于 SSE 流式输出
type StreamWriter struct {
	c *gin.Context
}

// WriteToken 输出 token
func (w *StreamWriter) WriteToken(token string) {
	w.WriteEvent(agent.EventToken, token)
}

// WriteEvent 输出事件
func (w *StreamWriter) WriteEvent(event agent.EventName, data string) {
	fmt.Fprintf(w.c.Writer, "event: %s\n", string(event))
	fmt.Fprintf(w.c.Writer, "data: %s\n\n", data)
	w.c.Writer.Flush()
}

// StreamResponse SSE 封装
func StreamResponse(c *gin.Context, fn func(writer *StreamWriter)) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	writer := &StreamWriter{c: c}
	fn(writer)
}

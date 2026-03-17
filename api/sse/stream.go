package httpapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// StreamWriter 用于 SSE 流式输出
type StreamWriter struct {
	c *gin.Context
}

// WriteToken 输出 token
func (w *StreamWriter) WriteToken(token string) {
	fmt.Fprintf(w.c.Writer, "event: token\n")
	fmt.Fprintf(w.c.Writer, "data: %s\n\n", token)
	w.c.Writer.Flush()
}

// WriteEvent 输出事件
func (w *StreamWriter) WriteEvent(event string, data string) {
	fmt.Fprintf(w.c.Writer, "event: %s\n", event)
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

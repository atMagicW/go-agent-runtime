package httpapi

import (
	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// SSEPublisher 将运行时事件发布到 SSE
type SSEPublisher struct {
	writer *StreamWriter
}

// NewSSEPublisher 创建 SSEPublisher
func NewSSEPublisher(writer *StreamWriter) *SSEPublisher {
	return &SSEPublisher{
		writer: writer,
	}
}

// Publish 发布事件
func (p *SSEPublisher) Publish(event agent.EventName, data string) {
	if p == nil || p.writer == nil {
		return
	}
	p.writer.WriteEvent(event, data)
}

var _ ports.EventPublisher = (*SSEPublisher)(nil)

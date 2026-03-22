package ports

import "github.com/atMagicW/go-agent-runtime/internal/domain/agent"

// EventPublisher 定义运行时事件发布接口
type EventPublisher interface {
	Publish(event agent.EventName, data string)
}

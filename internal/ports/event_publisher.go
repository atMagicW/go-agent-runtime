package ports

// EventPublisher 定义运行时事件发布接口
type EventPublisher interface {
	Publish(event string, data string)
}

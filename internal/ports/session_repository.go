package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// SessionRepository 定义会话持久化接口
type SessionRepository interface {
	// CreateSessionIfNotExists 如果会话不存在则创建
	CreateSessionIfNotExists(ctx context.Context, session agent.Session) error

	// GetSession 获取会话基本信息
	GetSession(ctx context.Context, sessionID string) (agent.Session, error)

	// SaveMessage 保存一条消息
	SaveMessage(ctx context.Context, sessionID string, message agent.Message) error

	// ListMessages 获取最近若干条消息，按时间正序返回
	ListMessages(ctx context.Context, sessionID string, limit int) ([]agent.Message, error)

	// SaveConversationState 保存会话状态
	SaveConversationState(ctx context.Context, state agent.ConversationState) error

	// GetConversationState 获取会话状态
	GetConversationState(ctx context.Context, sessionID string) (agent.ConversationState, error)
}

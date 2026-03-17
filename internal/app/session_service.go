package app

import (
	"context"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// SessionService 提供会话持久化能力
type SessionService struct {
	repo ports.SessionRepository
}

// NewSessionService 创建 SessionService
func NewSessionService(repo ports.SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

// EnsureSession 确保会话存在
func (s *SessionService) EnsureSession(ctx context.Context, sessionID string, userID string) error {
	return s.repo.CreateSessionIfNotExists(ctx, agent.Session{
		SessionID:    sessionID,
		UserID:       userID,
		Summary:      "",
		ActiveTaskID: "",
	})
}

// LoadConversationState 加载会话状态
func (s *SessionService) LoadConversationState(ctx context.Context, sessionID string) (agent.ConversationState, error) {
	return s.repo.GetConversationState(ctx, sessionID)
}

// SaveUserMessage 保存用户消息
func (s *SessionService) SaveUserMessage(ctx context.Context, sessionID string, content string) error {
	return s.repo.SaveMessage(ctx, sessionID, agent.Message{
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	})
}

// SaveAssistantMessage 保存助手消息
func (s *SessionService) SaveAssistantMessage(ctx context.Context, sessionID string, content string) error {
	return s.repo.SaveMessage(ctx, sessionID, agent.Message{
		Role:      "assistant",
		Content:   content,
		CreatedAt: time.Now(),
	})
}

// SaveConversationState 保存会话状态
func (s *SessionService) SaveConversationState(ctx context.Context, state agent.ConversationState) error {
	return s.repo.SaveConversationState(ctx, state)
}

// GetSession 获取会话基本信息
func (s *SessionService) GetSession(ctx context.Context, sessionID string) (agent.Session, error) {
	return s.repo.GetSession(ctx, sessionID)
}

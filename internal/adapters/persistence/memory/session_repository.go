package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

type SessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]agent.Session
	messages map[string][]agent.Message
	states   map[string]agent.ConversationState
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[string]agent.Session),
		messages: make(map[string][]agent.Message),
		states:   make(map[string]agent.ConversationState),
	}
}

func (r *SessionRepository) CreateSessionIfNotExists(ctx context.Context, session agent.Session) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[session.SessionID]; !ok {
		r.sessions[session.SessionID] = session
	}

	return nil
}

func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (agent.Session, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.sessions[sessionID]
	if !ok {
		return agent.Session{}, fmt.Errorf("session not found: %s", sessionID)
	}

	return s, nil
}

func (r *SessionRepository) SaveMessage(ctx context.Context, sessionID string, message agent.Message) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[sessionID] = append(r.messages[sessionID], message)
	return nil
}

func (r *SessionRepository) ListMessages(ctx context.Context, sessionID string, limit int) ([]agent.Message, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.messages[sessionID]
	if limit > 0 && len(items) > limit {
		items = items[len(items)-limit:]
	}

	out := make([]agent.Message, len(items))
	copy(out, items)
	return out, nil
}

func (r *SessionRepository) SaveConversationState(ctx context.Context, state agent.ConversationState) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	r.states[state.SessionID] = state
	return nil
}

func (r *SessionRepository) GetConversationState(ctx context.Context, sessionID string) (agent.ConversationState, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	if s, ok := r.states[sessionID]; ok {
		return s, nil
	}

	return agent.ConversationState{
		SessionID: sessionID,
		Messages:  r.messages[sessionID],
		Variables: map[string]any{},
	}, nil
}

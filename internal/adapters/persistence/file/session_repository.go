package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

type SessionRepository struct {
	mu       sync.Mutex
	dataDir  string
	sessions map[string]agent.Session
	messages map[string][]agent.Message
	states   map[string]agent.ConversationState
}

func NewSessionRepository(dataDir string) (*SessionRepository, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}

	r := &SessionRepository{
		dataDir:  dataDir,
		sessions: map[string]agent.Session{},
		messages: map[string][]agent.Message{},
		states:   map[string]agent.ConversationState{},
	}

	if err := r.load(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SessionRepository) CreateSessionIfNotExists(ctx context.Context, session agent.Session) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[session.SessionID]; !ok {
		r.sessions[session.SessionID] = session
		return r.save()
	}
	return nil
}

func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (agent.Session, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

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
	return r.save()
}

func (r *SessionRepository) ListMessages(ctx context.Context, sessionID string, limit int) ([]agent.Message, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

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
	return r.save()
}

func (r *SessionRepository) GetConversationState(ctx context.Context, sessionID string) (agent.ConversationState, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if s, ok := r.states[sessionID]; ok {
		return s, nil
	}

	return agent.ConversationState{
		SessionID: sessionID,
		Messages:  r.messages[sessionID],
		Variables: map[string]any{},
	}, nil
}

func (r *SessionRepository) load() error {
	if err := loadJSON(filepath.Join(r.dataDir, "sessions.json"), &r.sessions); err != nil {
		return err
	}
	if err := loadJSON(filepath.Join(r.dataDir, "messages.json"), &r.messages); err != nil {
		return err
	}
	if err := loadJSON(filepath.Join(r.dataDir, "states.json"), &r.states); err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) save() error {
	if err := saveJSON(filepath.Join(r.dataDir, "sessions.json"), r.sessions); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(r.dataDir, "messages.json"), r.messages); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(r.dataDir, "states.json"), r.states); err != nil {
		return err
	}
	return nil
}

func loadJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, target)
}

func saveJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

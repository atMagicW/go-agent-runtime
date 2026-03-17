package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// SessionRepository 是 PostgreSQL 的会话仓储实现
type SessionRepository struct {
	db *pgxpool.Pool
}

// NewSessionRepository 创建 PostgreSQL 会话仓储
func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// CreateSessionIfNotExists 如果会话不存在则创建
func (r *SessionRepository) CreateSessionIfNotExists(ctx context.Context, session agent.Session) error {
	const query = `
INSERT INTO sessions (session_id, user_id, summary, active_task_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (session_id) DO NOTHING
`
	_, err := r.db.Exec(
		ctx,
		query,
		session.SessionID,
		session.UserID,
		session.Summary,
		session.ActiveTaskID,
	)
	return err
}

// GetSession 获取会话
func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (agent.Session, error) {
	const query = `
SELECT session_id, user_id, summary, active_task_id, created_at, updated_at
FROM sessions
WHERE session_id = $1
`
	var s agent.Session
	err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&s.SessionID,
		&s.UserID,
		&s.Summary,
		&s.ActiveTaskID,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return agent.Session{}, err
	}
	return s, nil
}

// SaveMessage 保存消息
func (r *SessionRepository) SaveMessage(ctx context.Context, sessionID string, message agent.Message) error {
	const insertQuery = `
INSERT INTO session_messages (session_id, role, content, created_at)
VALUES ($1, $2, $3, $4)
`
	if _, err := r.db.Exec(ctx, insertQuery, sessionID, message.Role, message.Content, message.CreatedAt); err != nil {
		return err
	}

	const updateSessionQuery = `
UPDATE sessions
SET updated_at = NOW()
WHERE session_id = $1
`
	_, err := r.db.Exec(ctx, updateSessionQuery, sessionID)
	return err
}

// ListMessages 获取最近若干条消息，按时间正序返回
func (r *SessionRepository) ListMessages(ctx context.Context, sessionID string, limit int) ([]agent.Message, error) {
	const query = `
SELECT role, content, created_at
FROM (
    SELECT role, content, created_at
    FROM session_messages
    WHERE session_id = $1
    ORDER BY created_at DESC
    LIMIT $2
) t
ORDER BY created_at ASC
`
	rows, err := r.db.Query(ctx, query, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]agent.Message, 0)
	for rows.Next() {
		var m agent.Message
		if err := rows.Scan(&m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}

	return out, rows.Err()
}

// SaveConversationState 保存会话状态
func (r *SessionRepository) SaveConversationState(ctx context.Context, state agent.ConversationState) error {
	variablesJSON, err := json.Marshal(state.Variables)
	if err != nil {
		return err
	}

	const upsertStateQuery = `
INSERT INTO conversation_states (session_id, variables_json, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (session_id)
DO UPDATE SET
    variables_json = EXCLUDED.variables_json,
    updated_at = NOW()
`
	if _, err := r.db.Exec(ctx, upsertStateQuery, state.SessionID, variablesJSON); err != nil {
		return err
	}

	const updateSessionQuery = `
UPDATE sessions
SET summary = $2,
    active_task_id = $3,
    updated_at = NOW()
WHERE session_id = $1
`
	_, err = r.db.Exec(ctx, updateSessionQuery, state.SessionID, state.Summary, state.ActiveTaskID)
	return err
}

// GetConversationState 获取会话状态
func (r *SessionRepository) GetConversationState(ctx context.Context, sessionID string) (agent.ConversationState, error) {
	session, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return agent.ConversationState{}, err
	}

	messages, err := r.ListMessages(ctx, sessionID, 20)
	if err != nil {
		return agent.ConversationState{}, err
	}

	const query = `
SELECT variables_json
FROM conversation_states
WHERE session_id = $1
`
	var raw []byte
	err = r.db.QueryRow(ctx, query, sessionID).Scan(&raw)

	var variables map[string]any
	if err != nil {
		// 第一版：如果 conversation_states 还没初始化，不报错，返回空 variables
		variables = map[string]any{}
	} else {
		if len(raw) == 0 {
			variables = map[string]any{}
		} else {
			if unmarshalErr := json.Unmarshal(raw, &variables); unmarshalErr != nil {
				return agent.ConversationState{}, unmarshalErr
			}
		}
	}

	if variables == nil {
		variables = map[string]any{}
	}

	return agent.ConversationState{
		SessionID:    session.SessionID,
		Messages:     messages,
		Summary:      session.Summary,
		ActiveTaskID: session.ActiveTaskID,
		Variables:    variables,
	}, nil
}

// IsNotFoundError 第一版占位，后续如果要细化 pgx 错误可扩展
func IsNotFoundError(err error) bool {
	return err != nil && errors.Is(err, err)
}

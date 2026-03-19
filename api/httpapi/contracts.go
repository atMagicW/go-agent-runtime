package httpapi

import (
	"context"

	httpapi "github.com/atMagicW/go-agent-runtime/api/sse"
	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

// AgentService 定义 HTTP 层需要的 Agent 服务能力
type AgentService interface {
	Run(reqCtx agent.RequestContext, message string) (*agent.FinalResponse, error)
	RunStream(reqCtx agent.RequestContext, message string, writer *httpapi.StreamWriter)
}

// SessionService 定义 HTTP 层需要的 Session 服务能力
type SessionService interface {
	LoadConversationState(ctx context.Context, sessionID string) (agent.ConversationState, error)
	GetSession(ctx context.Context, sessionID string) (agent.Session, error)
}

// CapabilityService 定义能力列表查询接口
type CapabilityService interface {
	ListCapabilities() []capability.Descriptor
}

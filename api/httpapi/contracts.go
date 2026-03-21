package httpapi

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
	domainprompt "github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
)

// AgentService 定义 HTTP 层需要的 Agent 服务能力
type AgentService interface {
	Run(reqCtx agent.RequestContext, message string) (*agent.FinalResponse, error)
	RunStream(reqCtx agent.RequestContext, message string, writer *StreamWriter)
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

// IngestService 定义知识库写入接口
type IngestService interface {
	IngestText(ctx context.Context, req rag.IngestTextRequest) (*rag.IngestTextResponse, error)
}

// PromptService 定义 Prompt 查询接口
type PromptService interface {
	GetLatest(ctx context.Context, promptName string) (domainprompt.Template, error)
	ListVersions(ctx context.Context, promptName string) ([]domainprompt.Template, error)
}

// MCPService 定义 MCP server 查询接口
type MCPService interface {
	ListServers() []cfg.MCPServerConfig
}

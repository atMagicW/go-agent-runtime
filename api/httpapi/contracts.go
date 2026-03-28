package httpapi

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	domaincap "github.com/atMagicW/go-agent-runtime/internal/domain/capability"
	domainprompt "github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
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

type CapabilityView struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
	Version     string   `json:"version,omitempty"`
	Enabled     bool     `json:"enabled"`
	Source      string   `json:"source"`
	ServerName  string   `json:"server_name,omitempty"`
	RemoteTool  string   `json:"remote_tool,omitempty"`
}

// CapabilityService 定义能力列表查询接口
type CapabilityService interface {
	ListCapabilities() []CapabilityView
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

type MCPToolView struct {
	Name        string `json:"name"`
	RemoteTool  string `json:"remote_tool"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type MCPServerView struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Mode        string        `json:"mode"`
	BaseURL     string        `json:"base_url,omitempty"`
	ToolPath    string        `json:"tool_path,omitempty"`
	TimeoutMS   int           `json:"timeout_ms,omitempty"`
	Enabled     bool          `json:"enabled"`
	Tools       []MCPToolView `json:"tools"`
}

// MCPService 定义 MCP server 查询接口
type MCPService interface {
	ListServers() []MCPServerView
}

// SkillService 定义 Skill 查询接口
type SkillService interface {
	ListSkills() []domaincap.SkillDefinition
}

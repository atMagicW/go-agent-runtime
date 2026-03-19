package app

import (
	"context"
	"time"

	httpapi "github.com/atMagicW/go-agent-runtime/api/sse"
	memrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/memory"
	promptrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/prompt"
	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
	agentgov "github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
	agentintent "github.com/atMagicW/go-agent-runtime/internal/usecase/intent"
	agentplanner "github.com/atMagicW/go-agent-runtime/internal/usecase/planner"
	agentresponse "github.com/atMagicW/go-agent-runtime/internal/usecase/response"
	agentrouter "github.com/atMagicW/go-agent-runtime/internal/usecase/router"
	agentruntime "github.com/atMagicW/go-agent-runtime/internal/usecase/runtime"
)

// AgentService 是 Agent 运行时服务入口
type AgentService struct {
	orchestrator   *agentruntime.Orchestrator
	sessionService *SessionService
}

// capabilityRegistry 是本文件用到的最小能力注册接口
type capabilityRegistry interface {
	Get(name string) (ports.Capability, bool)
}

// NewAgentService 创建 AgentService
func NewAgentService(
	sessionService *SessionService,
	modelRouter *agentrouter.ModelRouter,
	registry capabilityRegistry,
) *AgentService {
	intentEngine := agentintent.NewEngine()
	planner := agentplanner.NewPlanner()

	breakers := agentgov.NewBreakerRegistry()
	fallbacks := agentgov.NewDefaultFallbackPolicy()

	capabilityRouter := agentrouter.NewCapabilityRouter(registry, breakers, fallbacks)
	ragRouter := agentrouter.NewRAGRouter(breakers, fallbacks)

	modelUsageRepo := memrepo.NewModelUsageRepository()
	auditRepo := memrepo.NewAuditRepository()
	promptRepository := promptrepo.NewInMemoryRepository()

	costTracker := agentgov.NewCostTracker(modelUsageRepo)
	auditLogger := agentgov.NewAuditLogger(auditRepo)
	responseComposer := agentresponse.NewTemplateResponseComposer(promptRepository)

	executor := agentruntime.NewPlanExecutor(
		modelRouter,
		capabilityRouter,
		ragRouter,
		responseComposer,
		costTracker,
	)

	orchestrator := agentruntime.NewOrchestrator(
		intentEngine,
		planner,
		executor,
		responseComposer,
		auditLogger,
	)

	return &AgentService{
		orchestrator:   orchestrator,
		sessionService: sessionService,
	}
}

// Run 非流式执行
func (s *AgentService) Run(
	reqCtx agent.RequestContext,
	message string,
) (*agent.FinalResponse, error) {
	baseCtx := context.Background()

	ctx, cancel := context.WithTimeout(baseCtx, 30*time.Second)
	defer cancel()

	// 1. 确保会话存在
	if err := s.sessionService.EnsureSession(ctx, reqCtx.SessionID, reqCtx.UserID); err != nil {
		return nil, err
	}

	// 2. 保存用户消息
	if err := s.sessionService.SaveUserMessage(ctx, reqCtx.SessionID, message); err != nil {
		return nil, err
	}

	// 3. 加载历史会话状态
	conversationState, err := s.sessionService.LoadConversationState(ctx, reqCtx.SessionID)
	if err != nil {
		return nil, err
	}

	runtimeCtx := agent.RuntimeContext{
		Request:      reqCtx,
		Conversation: conversationState,
		Variables:    conversationState.Variables,
	}
	if runtimeCtx.Variables == nil {
		runtimeCtx.Variables = map[string]any{}
	}

	// 4. 执行主链路
	resp, _, err := s.orchestrator.Run(ctx, runtimeCtx, message)
	if err != nil {
		return nil, err
	}

	// 5. 保存助手回复
	if err := s.sessionService.SaveAssistantMessage(ctx, reqCtx.SessionID, resp.Message); err != nil {
		return nil, err
	}

	// 6. 保存最新会话状态
	conversationState.Variables = runtimeCtx.Variables

	if err := s.sessionService.SaveConversationState(ctx, conversationState); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RunStream 流式执行
func (s *AgentService) RunStream(
	reqCtx agent.RequestContext,
	message string,
	writer *httpapi.StreamWriter,
) {
	baseCtx := context.Background()

	ctx, cancel := context.WithTimeout(baseCtx, 30*time.Second)
	defer cancel()

	if err := s.sessionService.EnsureSession(ctx, reqCtx.SessionID, reqCtx.UserID); err != nil {
		writer.WriteEvent("error", err.Error())
		return
	}

	if err := s.sessionService.SaveUserMessage(ctx, reqCtx.SessionID, message); err != nil {
		writer.WriteEvent("error", err.Error())
		return
	}

	conversationState, err := s.sessionService.LoadConversationState(ctx, reqCtx.SessionID)
	if err != nil {
		writer.WriteEvent("error", err.Error())
		return
	}

	runtimeCtx := agent.RuntimeContext{
		Request:      reqCtx,
		Conversation: conversationState,
		Variables:    conversationState.Variables,
	}
	if runtimeCtx.Variables == nil {
		runtimeCtx.Variables = map[string]any{}
	}

	writer.WriteEvent("plan", "starting orchestrator")

	resp, results, err := s.orchestrator.Run(ctx, runtimeCtx, message)
	if err != nil {
		writer.WriteEvent("error", err.Error())
		return
	}

	for _, result := range results {
		if result.Success {
			writer.WriteEvent("step_done", result.StepID)
		} else {
			writer.WriteEvent("step_error", result.StepID+":"+result.Error)
		}
	}

	if err := s.sessionService.SaveAssistantMessage(ctx, reqCtx.SessionID, resp.Message); err != nil {
		writer.WriteEvent("error", err.Error())
		return
	}

	conversationState.Variables = runtimeCtx.Variables

	if err := s.sessionService.SaveConversationState(ctx, conversationState); err != nil {
		writer.WriteEvent("error", err.Error())
		return
	}

	writer.WriteToken(resp.Message)
	writer.WriteEvent("done", "completed")
}

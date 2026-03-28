package app

import (
	"context"
	"fmt"
	"time"

	"github.com/atMagicW/go-agent-runtime/api/httpapi"
	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
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
	orchestrator     *agentruntime.Orchestrator
	sessionService   *SessionService
	intentEngine     ports.IntentEngine
	modelRouter      ports.ModelRouter
	responseComposer ports.ResponseComposer
}

// capabilityRegistry 是本文件用到的最小能力注册接口
type capabilityRegistry interface {
	Get(name string) (ports.Capability, bool)
}
type ragSearchService interface {
	Search(ctx context.Context, kbID string, query string, topK int) ([]rag.Evidence, error)
}

// NewAgentService 创建 AgentService
func NewAgentService(
	sessionService *SessionService,
	modelRouter *agentrouter.ModelRouter,
	registry capabilityRegistry,
	ragService ragSearchService,
	promptRepository ports.PromptRepository,
	modelUsageRepo ports.ModelUsageRepository,
	auditRepo ports.AuditRepository,
	breakers *agentgov.BreakerRegistry,
	fallbacks *agentgov.FallbackPolicy,
) *AgentService {
	ruleClassifier := agentintent.NewRuleClassifier()
	llmClassifier := agentintent.NewLLMClassifier(modelRouter)
	intentEngine := agentintent.NewEngine(
		ruleClassifier,
		llmClassifier,
	)
	planner := agentplanner.NewPlanner()

	capabilityRouter := agentrouter.NewCapabilityRouter(registry, breakers, fallbacks)
	ragRouter := agentrouter.NewRAGRouter(ragService, breakers, fallbacks)

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
		orchestrator:     orchestrator,
		sessionService:   sessionService,
		intentEngine:     intentEngine,
		modelRouter:      modelRouter,
		responseComposer: responseComposer,
	}
}

// Run 非流式执行
func (s *AgentService) Run(
	reqCtx agent.RequestContext,
	message string,
) (*agent.FinalResponse, error) {
	baseCtx := context.Background()

	ctx, cancel := context.WithTimeout(baseCtx, time.Duration(agent.DefaultRequestTimeoutMS)*time.Millisecond)
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
		writer.WriteEvent(agent.EventError, err.Error())
		return
	}

	if err := s.sessionService.SaveUserMessage(ctx, reqCtx.SessionID, message); err != nil {
		writer.WriteEvent(agent.EventError, err.Error())
		return
	}

	conversationState, err := s.sessionService.LoadConversationState(ctx, reqCtx.SessionID)
	if err != nil {
		writer.WriteEvent(agent.EventError, err.Error())
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

	publisher := httpapi.NewSSEPublisher(writer)

	intentResult, err := s.orchestratorIntentOnly(ctx, runtimeCtx, message)
	if err != nil {
		writer.WriteEvent(agent.EventError, err.Error())
		return
	}

	// chat / write 仍然优先直接走模型 token streaming
	if intentResult.IntentType == agent.IntentChat || intentResult.IntentType == agent.IntentWrite {
		prompt := "请回答用户请求：\n" + message

		err := s.streamDirectModel(ctx, runtimeCtx, prompt, writer)
		if err != nil {
			writer.WriteEvent(agent.EventError, err.Error())
			return
		}

		if err := s.sessionService.SaveAssistantMessage(ctx, reqCtx.SessionID, "[streamed response]"); err != nil {
			writer.WriteEvent(agent.EventError, err.Error())
			return
		}

		writer.WriteEvent(agent.EventDone, "completed")
		return
	}

	// workflow / rag / tool 使用事件版执行
	resp, results, err := s.orchestrator.RunWithEvents(ctx, runtimeCtx, message, publisher)
	if err != nil {
		writer.WriteEvent(agent.EventError, err.Error())
		return
	}

	// 先发一个标记，前端可切换成“最终回答中”
	writer.WriteEvent(agent.EventFinalAnswerStart, "streaming")

	if err := s.streamWorkflowFinalAnswer(ctx, runtimeCtx, message, results, writer); err != nil {
		// 如果流式最终回答失败，回退为一次性输出
		writer.WriteEvent(agent.EventFinalAnswerFallback, resp.Message)
		writer.WriteToken(resp.Message)
	}

	if err := s.sessionService.SaveAssistantMessage(ctx, reqCtx.SessionID, resp.Message); err != nil {
		writer.WriteEvent(agent.EventError, err.Error())
		return
	}

	conversationState.Variables = runtimeCtx.Variables
	if err := s.sessionService.SaveConversationState(ctx, conversationState); err != nil {
		writer.WriteEvent(agent.EventError, err.Error())
		return
	}

	writer.WriteEvent(agent.EventDone, "completed")
}

func (s *AgentService) orchestratorIntentOnly(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	message string,
) (agent.IntentResult, error) {
	return s.intentEngine.Recognize(ctx, runtimeCtx, message)
}

func (s *AgentService) streamDirectModel(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	prompt string,
	writer *httpapi.StreamWriter,
) error {
	if s.modelRouter == nil {
		return fmt.Errorf("model router is nil")
	}

	return s.modelRouter.GenerateStream(ctx, runtimeCtx, ports.ModelCallRequest{
		TaskType: "llm_generate",
		Prompt:   prompt,
		Stream:   true,
	}, func(text string) error {
		writer.WriteToken(text)
		return nil
	})
}

func (s *AgentService) streamWorkflowFinalAnswer(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	message string,
	results []agent.StepResult,
	writer *httpapi.StreamWriter,
) error {
	if s.modelRouter == nil || s.responseComposer == nil {
		return fmt.Errorf("model router or response composer is nil")
	}

	prompt, err := s.responseComposer.BuildPrompt(ctx, runtimeCtx, ports.ComposeRequest{
		Message:     message,
		PromptName:  string(prompt.PromptFinalResponse),
		StepResults: results,
	})
	if err != nil {
		return err
	}

	return s.modelRouter.GenerateStream(ctx, runtimeCtx, ports.ModelCallRequest{
		TaskType: "retrieve_answer",
		Prompt:   prompt,
		Stream:   true,
	}, func(text string) error {
		writer.WriteToken(text)
		return nil
	})
}

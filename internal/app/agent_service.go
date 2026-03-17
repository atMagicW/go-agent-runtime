package app

import (
	"context"
	"time"

	"github.com/atMagicW/go-agent-runtime/api/sse"
	memrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/persistence/memory"
	promptrepo "github.com/atMagicW/go-agent-runtime/internal/adapters/prompt"
	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	agentgov "github.com/atMagicW/go-agent-runtime/internal/usecase/governance"
	agentintent "github.com/atMagicW/go-agent-runtime/internal/usecase/intent"
	agentplanner "github.com/atMagicW/go-agent-runtime/internal/usecase/planner"
	agentrouter "github.com/atMagicW/go-agent-runtime/internal/usecase/router"
	agentruntime "github.com/atMagicW/go-agent-runtime/internal/usecase/runtime"
)

// AgentService 是 Agent 运行时服务入口
type AgentService struct {
	orchestrator  *agentruntime.Orchestrator
	promptService *PromptService
}

// NewAgentService 创建 AgentService
func NewAgentService() *AgentService {
	intentEngine := agentintent.NewEngine()
	planner := agentplanner.NewPlanner()

	modelRouter := agentrouter.NewModelRouter()
	capabilityRouter := agentrouter.NewCapabilityRouter()
	ragRouter := agentrouter.NewRAGRouter()

	modelUsageRepo := memrepo.NewModelUsageRepository()
	auditRepo := memrepo.NewAuditRepository()
	promptRepository := promptrepo.NewInMemoryRepository()

	_ = promptRepository // 先初始化，下一轮会接到 response composer

	costTracker := agentgov.NewCostTracker(modelUsageRepo)
	auditLogger := agentgov.NewAuditLogger(auditRepo)

	executor := agentruntime.NewPlanExecutor(
		modelRouter,
		capabilityRouter,
		ragRouter,
		costTracker,
	)

	orchestrator := agentruntime.NewOrchestrator(
		intentEngine,
		planner,
		executor,
		auditLogger,
	)

	return &AgentService{
		orchestrator: orchestrator,
	}
}

// Run 非流式执行
func (s *AgentService) Run(
	reqCtx agent.RequestContext,
	message string,
) (*agent.FinalResponse, error) {
	baseCtx := context.Background()
	// TODO:
	// 1. 构建上下文
	// 2. 意图识别
	// 3. 生成计划
	// 4. 执行计划
	// 5. 返回最终回复

	ctx, cancel := context.WithTimeout(baseCtx, 30*time.Second)
	defer cancel()

	runtimeCtx := agent.RuntimeContext{
		Request: reqCtx,
		Conversation: agent.ConversationState{
			SessionID: reqCtx.SessionID,
		},
		Variables: map[string]any{},
	}

	resp, _, err := s.orchestrator.Run(ctx, runtimeCtx, message)
	if err != nil {
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

	runtimeCtx := agent.RuntimeContext{
		Request: reqCtx,
		Conversation: agent.ConversationState{
			SessionID: reqCtx.SessionID,
		},
		Variables: map[string]any{},
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

	writer.WriteToken(resp.Message)
	writer.WriteEvent("done", "completed")
}

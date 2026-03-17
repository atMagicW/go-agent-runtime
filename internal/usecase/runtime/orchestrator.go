package runtime

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// AuditLogger 定义审计记录接口，避免 usecase 之间强耦合
type AuditLogger interface {
	Log(
		ctx context.Context,
		reqCtx agent.RequestContext,
		intent agent.IntentResult,
		plan agent.ExecutionPlan,
		results []agent.StepResult,
		finalResp agent.FinalResponse,
		startedAt time.Time,
		status string,
	) error
}

// Orchestrator 是 Agent 运行时编排器
type Orchestrator struct {
	intentEngine ports.IntentEngine
	planner      ports.Planner
	executor     ports.Executor

	responseComposer ports.ResponseComposer
	auditLogger      AuditLogger
}

// NewOrchestrator 创建编排器
func NewOrchestrator(
	intentEngine ports.IntentEngine,
	planner ports.Planner,
	executor ports.Executor,
	responseComposer ports.ResponseComposer,
	auditLogger AuditLogger,
) *Orchestrator {
	return &Orchestrator{
		intentEngine:     intentEngine,
		planner:          planner,
		executor:         executor,
		responseComposer: responseComposer,
		auditLogger:      auditLogger,
	}
}

// Run 执行主链路
func (o *Orchestrator) Run(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	message string,
) (agent.FinalResponse, []agent.StepResult, error) {
	startedAt := time.Now()

	intentResult, err := o.intentEngine.Recognize(ctx, runtimeCtx, message)
	if err != nil {
		return agent.FinalResponse{}, nil, err
	}

	runtimeCtx.Intent = intentResult

	plan, err := o.planner.BuildPlan(ctx, runtimeCtx, message)
	if err != nil {
		return agent.FinalResponse{}, nil, err
	}

	results, err := o.executor.ExecutePlan(ctx, runtimeCtx, plan)
	if err != nil {
		finalResp := agent.FinalResponse{
			Message: "execution failed",
		}

		if o.auditLogger != nil {
			_ = o.auditLogger.Log(
				ctx,
				runtimeCtx.Request,
				intentResult,
				plan,
				results,
				finalResp,
				startedAt,
				"failed",
			)
		}

		return finalResp, results, err
	}

	totalCost := 0.0
	totalTokens := 0

	for _, result := range results {
		if !result.Success {
			continue
		}
		if cost, ok := result.Output["cost"].(float64); ok {
			totalCost += cost
		}
		if tokens, ok := result.Output["tokens"].(int); ok {
			totalTokens += tokens
		}
	}

	finalResp := agent.FinalResponse{
		Message: "execution completed",
		Cost:    totalCost,
		Tokens:  totalTokens,
	}

	if o.responseComposer != nil {
		composed, composeErr := o.responseComposer.Compose(
			ctx,
			runtimeCtx,
			ports.ComposeRequest{
				Message:     message,
				PromptName:  "final_response",
				StepResults: results,
			},
		)
		if composeErr == nil {
			finalResp.Message = composed.Text
			finalResp.Cost += composed.Cost
			finalResp.Tokens += composed.Tokens
		}
	}

	if o.auditLogger != nil {
		_ = o.auditLogger.Log(
			ctx,
			runtimeCtx.Request,
			intentResult,
			plan,
			results,
			finalResp,
			startedAt,
			"succeeded",
		)
	}

	return finalResp, results, nil
}

// ------------------------------------------------------------------

// CostTracker 定义成本跟踪接口
type CostTracker interface {
	Track(
		ctx context.Context,
		requestID string,
		sessionID string,
		provider string,
		modelName string,
		promptTokens int,
		completionTokens int,
		cost float64,
		latencyMS int64,
	) error
}

// PlanExecutor 执行器
type PlanExecutor struct {
	modelRouter      ports.ModelRouter
	capabilityRouter ports.CapabilityRouter
	ragRouter        ports.RAGRouter
	responseComposer ports.ResponseComposer

	costTracker CostTracker
}

// NewPlanExecutor 创建执行器
func NewPlanExecutor(
	modelRouter ports.ModelRouter,
	capabilityRouter ports.CapabilityRouter,
	ragRouter ports.RAGRouter,
	responseComposer ports.ResponseComposer,
	costTracker CostTracker,
) *PlanExecutor {
	return &PlanExecutor{
		modelRouter:      modelRouter,
		capabilityRouter: capabilityRouter,
		ragRouter:        ragRouter,
		responseComposer: responseComposer,
		costTracker:      costTracker,
	}
}

// ExecutePlan 执行计划
func (e *PlanExecutor) ExecutePlan(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	plan agent.ExecutionPlan,
) ([]agent.StepResult, error) {
	stepMap := make(map[string]agent.PlanStep, len(plan.Steps))
	for _, step := range plan.Steps {
		stepMap[step.StepID] = step
	}

	completed := make(map[string]agent.StepResult)
	results := make([]agent.StepResult, 0, len(plan.Steps))

	for len(completed) < len(plan.Steps) {
		ready := e.findReadySteps(plan.Steps, completed)

		if len(ready) == 0 {
			return results, fmt.Errorf("no executable step found, maybe dependency cycle")
		}

		grouped := e.groupByParallelKey(ready)

		for _, steps := range grouped {
			// 单个步骤直接执行；多个步骤并发执行
			if len(steps) == 1 {
				step := steps[0]
				result := e.executeStepWithRetry(ctx, runtimeCtx, step, completed)
				completed[step.StepID] = result
				results = append(results, result)
				continue
			}

			var wg sync.WaitGroup
			var mu sync.Mutex

			tmpResults := make([]agent.StepResult, 0, len(steps))

			for _, step := range steps {
				wg.Add(1)

				go func(s agent.PlanStep) {
					defer wg.Done()

					result := e.executeStepWithRetry(ctx, runtimeCtx, s, completed)

					mu.Lock()
					defer mu.Unlock()
					tmpResults = append(tmpResults, result)
				}(step)
			}

			wg.Wait()

			for _, result := range tmpResults {
				completed[result.StepID] = result
				results = append(results, result)
			}
		}
	}

	return results, nil
}

// findReadySteps 找到当前可执行步骤
func (e *PlanExecutor) findReadySteps(
	steps []agent.PlanStep,
	completed map[string]agent.StepResult,
) []agent.PlanStep {
	ready := make([]agent.PlanStep, 0)

	for _, step := range steps {
		if _, ok := completed[step.StepID]; ok {
			continue
		}

		allDepsDone := true
		for _, dep := range step.DependsOn {
			if _, ok := completed[dep]; !ok {
				allDepsDone = false
				break
			}
		}

		if allDepsDone {
			ready = append(ready, step)
		}
	}

	return ready
}

// groupByParallelKey 按并行组聚合
func (e *PlanExecutor) groupByParallelKey(steps []agent.PlanStep) map[string][]agent.PlanStep {
	grouped := make(map[string][]agent.PlanStep)

	for _, step := range steps {
		key := step.ParallelGroup
		if key == "" {
			key = "single:" + step.StepID
		}
		grouped[key] = append(grouped[key], step)
	}

	return grouped
}

// executeStepWithRetry 带重试执行步骤
func (e *PlanExecutor) executeStepWithRetry(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	step agent.PlanStep,
	completed map[string]agent.StepResult,
) agent.StepResult {
	maxRetries := step.RetryPolicy.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}

	var lastResult agent.StepResult
	for attempt := 0; attempt <= maxRetries; attempt++ {
		lastResult = e.executeStep(ctx, runtimeCtx, step, completed)
		if lastResult.Success {
			return lastResult
		}
		time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
	}

	return lastResult
}

// executeStep 执行单个步骤
func (e *PlanExecutor) executeStep(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	step agent.PlanStep,
	completed map[string]agent.StepResult,
) agent.StepResult {
	start := time.Now()

	stepCtx := ctx
	if step.TimeoutMS > 0 {
		var cancel context.CancelFunc
		stepCtx, cancel = context.WithTimeout(ctx, time.Duration(step.TimeoutMS)*time.Millisecond)
		defer cancel()
	}

	result := agent.StepResult{
		StepID:    step.StepID,
		StartedAt: start,
		EndedAt:   time.Now(),
	}

	switch step.Executor {
	case "model_router":
		resp, err := e.modelRouter.Generate(stepCtx, runtimeCtx, ports.ModelCallRequest{
			TaskType: string(step.Type),
			Model:    runtimeCtx.Request.Model,
			Prompt:   e.buildModelPrompt(step, completed),
			Stream:   runtimeCtx.Request.Stream,
		})
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.EndedAt = time.Now()
			return result
		}

		result.Success = true
		result.Output = map[string]any{
			"text":     resp.Text,
			"tokens":   resp.Tokens,
			"cost":     resp.Cost,
			"model":    resp.Model,
			"provider": resp.Provider,
		}
		result.EndedAt = time.Now()
		if e.costTracker != nil {
			_ = e.costTracker.Track(
				stepCtx,
				runtimeCtx.Request.RequestID,
				runtimeCtx.Request.SessionID,
				resp.Provider,
				resp.Model,
				64, // 第一版先写 mock prompt tokens
				resp.Tokens-64,
				resp.Cost,
				time.Since(start).Milliseconds(),
			)
		}
		return result

	case "capability_router":
		name, _ := step.Input["name"].(string)

		resp, err := e.capabilityRouter.Invoke(stepCtx, runtimeCtx, ports.CapabilityCallRequest{
			Name:  name,
			Input: step.Input,
		})
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.EndedAt = time.Now()
			return result
		}

		result.Success = true
		result.Output = resp.Output
		result.EndedAt = time.Now()
		return result

	case "rag_router":
		query, _ := step.Input["query"].(string)
		kb, _ := step.Input["kb"].(string)

		topK := 5
		if rawTopK, ok := step.Input["top_k"].(int); ok {
			topK = rawTopK
		}

		resp, err := e.ragRouter.Retrieve(stepCtx, runtimeCtx, ports.RetrievalRequest{
			KnowledgeBase: kb,
			Query:         query,
			TopK:          topK,
		})
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.EndedAt = time.Now()
			return result
		}

		result.Success = true
		result.Output = map[string]any{
			"knowledge_base": kb,
			"evidences":      resp.Evidences,
		}
		result.EndedAt = time.Now()
		return result
	case "response_composer":
		message, _ := step.Input["message"].(string)

		resp, err := e.responseComposer.Compose(stepCtx, runtimeCtx, ports.ComposeRequest{
			Message:     message,
			PromptName:  "final_response",
			StepResults: flattenCompletedResults(completed),
		})
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.EndedAt = time.Now()
			return result
		}

		result.Success = true
		result.Output = map[string]any{
			"text":   resp.Text,
			"tokens": resp.Tokens,
			"cost":   resp.Cost,
			"model":  resp.Model,
		}
		result.EndedAt = time.Now()
		return result
	default:
		result.Success = false
		result.Error = "unknown executor: " + step.Executor
		result.EndedAt = time.Now()
		return result
	}
}

// buildModelPrompt 构建给模型的 prompt
func (e *PlanExecutor) buildModelPrompt(
	step agent.PlanStep,
	completed map[string]agent.StepResult,
) string {
	prompt := "请基于当前步骤生成结果。\n"

	if msg, ok := step.Input["message"].(string); ok {
		prompt += "用户请求：" + msg + "\n"
	}

	for stepID, result := range completed {
		if !result.Success {
			continue
		}
		prompt += fmt.Sprintf("前序步骤[%s]输出：%v\n", stepID, result.Output)
	}

	return prompt
}

// flattenCompletedResults 将 completed map 转成切片
func flattenCompletedResults(completed map[string]agent.StepResult) []agent.StepResult {
	out := make([]agent.StepResult, 0, len(completed))
	for _, r := range completed {
		out = append(out, r)
	}
	return out
}

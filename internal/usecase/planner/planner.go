package planner

import (
	"context"
	"strings"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
	"github.com/google/uuid"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// Planner 是第一版任务规划器
type Planner struct {
}

// NewPlanner 创建 Planner
func NewPlanner() *Planner {
	return &Planner{}
}

// BuildPlan 根据意图生成执行计划
func (p *Planner) BuildPlan(
	_ context.Context,
	runtimeCtx agent.RuntimeContext,
	message string,
) (agent.ExecutionPlan, error) {
	plan := agent.ExecutionPlan{
		PlanID:    uuid.New().String(),
		Goal:      message,
		Status:    agent.PlanStatusPending,
		CreatedAt: time.Now(),
	}

	switch runtimeCtx.Intent.IntentType {
	case agent.IntentRetrievalQA:
		plan.Steps = []agent.PlanStep{
			{
				StepID:      "step_retrieve",
				Name:        "检索知识库",
				Type:        agent.StepTypeRetrieve,
				Executor:    agent.ExecutorRAGRouter,
				TimeoutMS:   agent.DefaultStepTimeoutMS,
				RetryPolicy: agent.RetryPolicy{MaxRetries: 1},
				Status:      agent.StepStatusPending,
				Input: map[string]any{
					"query": message,
					"kb":    "default",
					"top_k": 5,
				},
			},
			{
				StepID:      "step_compose",
				Name:        "生成最终回答",
				Type:        agent.StepTypeComposeResponse,
				Executor:    agent.ExecutorResponseComposer,
				DependsOn:   []string{"step_retrieve"},
				TimeoutMS:   agent.DefaultComposeTimeoutMS,
				RetryPolicy: agent.RetryPolicy{MaxRetries: 1},
				Status:      agent.StepStatusPending,
				Input: map[string]any{
					"message": message,
				},
			},
		}

	case agent.IntentToolCall:
		plan.Steps = []agent.PlanStep{
			{
				StepID:      "step_tool",
				Name:        "调用能力",
				Type:        agent.StepTypeTool,
				Executor:    agent.ExecutorCapabilityRouter,
				TimeoutMS:   agent.DefaultStepTimeoutMS,
				RetryPolicy: agent.RetryPolicy{MaxRetries: 2},
				Status:      agent.StepStatusPending,
				Input: map[string]any{
					"name":    pickCapabilityName(message),
					"message": message,
				},
			},
			{
				StepID:      "step_compose",
				Name:        "生成最终回答",
				Type:        agent.StepTypeComposeResponse,
				Executor:    agent.ExecutorResponseComposer,
				DependsOn:   []string{"step_tool"},
				TimeoutMS:   agent.DefaultComposeTimeoutMS,
				RetryPolicy: agent.RetryPolicy{MaxRetries: 1},
				Status:      agent.StepStatusPending,
				Input: map[string]any{
					"message": message,
				},
			},
		}

	case agent.IntentWorkflow:
		plan.Steps = []agent.PlanStep{
			{
				StepID:        "step_retrieve_1",
				Name:          "并发检索知识库A",
				Type:          agent.StepTypeRetrieve,
				Executor:      agent.ExecutorRAGRouter,
				ParallelGroup: "group_retrieve",
				TimeoutMS:     agent.DefaultRagTimeoutMS,
				RetryPolicy:   agent.RetryPolicy{MaxRetries: 1},
				Status:        agent.StepStatusPending,
				Input: map[string]any{
					"query": message,
					"kb":    "knowledge_a",
					"top_k": 3,
				},
			},
			{
				StepID:        "step_retrieve_2",
				Name:          "并发检索知识库B",
				Type:          agent.StepTypeRetrieve,
				Executor:      agent.ExecutorRAGRouter,
				ParallelGroup: "group_retrieve",
				TimeoutMS:     agent.DefaultRagTimeoutMS,
				RetryPolicy:   agent.RetryPolicy{MaxRetries: 1},
				Status:        agent.StepStatusPending,
				Input: map[string]any{
					"query": message,
					"kb":    "knowledge_b",
					"top_k": 3,
				},
			},
			{
				StepID:      "step_tool",
				Name:        "调用本地能力",
				Type:        agent.StepTypeTool,
				Executor:    agent.ExecutorCapabilityRouter,
				DependsOn:   []string{"step_retrieve_1", "step_retrieve_2"},
				TimeoutMS:   agent.DefaultStepTimeoutMS,
				RetryPolicy: agent.RetryPolicy{MaxRetries: 2},
				Status:      agent.StepStatusPending,
				Input: map[string]any{
					"name":    pickCapabilityName(message),
					"message": message,
				},
			},
			{
				StepID:      "step_compose",
				Name:        "汇总生成结果",
				Type:        agent.StepTypeComposeResponse,
				Executor:    agent.ExecutorResponseComposer,
				DependsOn:   []string{"step_tool"},
				TimeoutMS:   agent.DefaultLongStepTimeoutMS,
				RetryPolicy: agent.RetryPolicy{MaxRetries: 1},
				Status:      agent.StepStatusPending,
				Input: map[string]any{
					"message": message,
				},
			},
		}

	default:
		plan.Steps = []agent.PlanStep{
			{
				StepID:      "step_chat",
				Name:        "直接生成回答",
				Type:        agent.StepTypeLLMGenerate,
				Executor:    agent.ExecutorResponseComposer,
				TimeoutMS:   agent.DefaultComposeTimeoutMS,
				RetryPolicy: agent.RetryPolicy{MaxRetries: 1},
				Status:      agent.StepStatusPending,
				Input: map[string]any{
					"message": message,
				},
			},
		}
	}

	return plan, nil
}

func pickCapabilityName(message string) string {
	lower := strings.ToLower(message)

	switch {
	case strings.Contains(lower, "新闻"), strings.Contains(lower, "news"):
		return capability.CapabilityMCPNewsSearch
	case strings.Contains(lower, "文档"), strings.Contains(lower, "doc"):
		return capability.CapabilityMCPDocLookup
	case strings.Contains(lower, "mcp"), strings.Contains(lower, "远程"), strings.Contains(lower, "搜索"):
		return capability.CapabilityMCPWebSearch
	case strings.Contains(lower, "关键词"), strings.Contains(lower, "keyword"):
		return capability.CapabilityKeywordExtract
	default:
		return capability.CapabilityResumeAnalyzer
	}
}

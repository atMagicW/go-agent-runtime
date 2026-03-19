package planner

import (
	"context"
	"strings"
	"time"

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
				Executor:    "rag_router",
				TimeoutMS:   5000,
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
				Executor:    "response_composer",
				DependsOn:   []string{"step_retrieve"},
				TimeoutMS:   8000,
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
				Executor:    "capability_router",
				TimeoutMS:   5000,
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
				Executor:    "response_composer",
				DependsOn:   []string{"step_tool"},
				TimeoutMS:   8000,
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
				Executor:      "rag_router",
				ParallelGroup: "group_retrieve",
				TimeoutMS:     4000,
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
				Executor:      "rag_router",
				ParallelGroup: "group_retrieve",
				TimeoutMS:     4000,
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
				Executor:    "capability_router",
				DependsOn:   []string{"step_retrieve_1", "step_retrieve_2"},
				TimeoutMS:   5000,
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
				Executor:    "response_composer",
				DependsOn:   []string{"step_tool"},
				TimeoutMS:   10000,
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
				Executor:    "response_composer",
				TimeoutMS:   8000,
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
	case strings.Contains(lower, "mcp"), strings.Contains(lower, "远程"), strings.Contains(lower, "搜索"):
		return "mcp_web_search"
	case strings.Contains(lower, "关键词"), strings.Contains(lower, "keyword"):
		return "keyword_extract_tool"
	default:
		return "resume_analyzer"
	}
}

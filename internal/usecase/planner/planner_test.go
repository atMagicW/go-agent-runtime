package planner

import (
	"context"
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

func TestBuildPlan_Retrieval(t *testing.T) {
	p := NewPlanner()

	runtimeCtx := agent.RuntimeContext{
		Intent: agent.IntentResult{
			IntentType: agent.IntentRetrievalQA,
		},
	}

	plan, err := p.BuildPlan(context.Background(), runtimeCtx, "请检索知识库")
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}

	if len(plan.Steps) != 2 {
		t.Fatalf("BuildPlan() steps len = %d, want 2", len(plan.Steps))
	}

	if plan.Steps[0].Executor != "rag_router" {
		t.Fatalf("first step executor = %s", plan.Steps[0].Executor)
	}

	if plan.Steps[1].Executor != "response_composer" {
		t.Fatalf("second step executor = %s", plan.Steps[1].Executor)
	}
}

package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// Executor 定义计划执行器
type Executor interface {
	ExecutePlan(ctx context.Context, runtimeCtx agent.RuntimeContext, plan agent.ExecutionPlan) ([]agent.StepResult, error)
}

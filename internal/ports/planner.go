package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
)

// Planner 定义任务规划器
type Planner interface {
	BuildPlan(ctx context.Context, runtimeCtx agent.RuntimeContext, message string) (agent.ExecutionPlan, error)
}

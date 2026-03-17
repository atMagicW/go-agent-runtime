package governance

import (
	"context"
	"time"

	"github.com/atMagicW/go-agent-runtime/internal/domain/model"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// CostTracker 负责记录模型调用成本
type CostTracker struct {
	repo ports.ModelUsageRepository
}

// NewCostTracker 创建成本统计器
func NewCostTracker(repo ports.ModelUsageRepository) *CostTracker {
	return &CostTracker{
		repo: repo,
	}
}

// Track 记录一次模型调用
func (t *CostTracker) Track(
	ctx context.Context,
	requestID string,
	sessionID string,
	provider string,
	modelName string,
	promptTokens int,
	completionTokens int,
	cost float64,
	latencyMS int64,
) error {
	record := model.UsageRecord{
		RequestID:        requestID,
		SessionID:        sessionID,
		Provider:         provider,
		Model:            modelName,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		Cost:             cost,
		LatencyMS:        latencyMS,
		CreatedAt:        time.Now(),
	}

	return t.repo.Save(ctx, record)
}

package governance

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/domain/audit"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// AuditLogger 负责记录请求级审计日志
type AuditLogger struct {
	repo ports.AuditRepository
}

// NewAuditLogger 创建审计日志器
func NewAuditLogger(repo ports.AuditRepository) *AuditLogger {
	return &AuditLogger{
		repo: repo,
	}
}

// Log 记录一次完整执行的审计信息
func (l *AuditLogger) Log(
	ctx context.Context,
	reqCtx agent.RequestContext,
	intent agent.IntentResult,
	plan agent.ExecutionPlan,
	results []agent.StepResult,
	finalResp agent.FinalResponse,
	startedAt time.Time,
	status string,
) error {
	planBytes, _ := json.Marshal(plan)

	models := collectModels(results)
	kbs := collectKnowledgeBases(results)
	caps := collectCapabilities(results)

	record := audit.Record{
		AuditID:        uuid.New().String(),
		RequestID:      reqCtx.RequestID,
		SessionID:      reqCtx.SessionID,
		UserID:         reqCtx.UserID,
		Intent:         string(intent.IntentType),
		PlanSnapshot:   string(planBytes),
		ModelsUsed:     models,
		KnowledgeBases: kbs,
		Capabilities:   caps,
		ResultSummary:  finalResp.Message,
		Status:         status,
		DurationMS:     time.Since(startedAt).Milliseconds(),
		TotalCost:      finalResp.Cost,
		CreatedAt:      time.Now(),
	}

	return l.repo.Save(ctx, record)
}

// collectModels 收集执行中用到的模型
func collectModels(results []agent.StepResult) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0)

	for _, r := range results {
		if !r.Success {
			continue
		}
		if modelName, ok := r.Output["model"].(string); ok && modelName != "" {
			if _, exists := seen[modelName]; !exists {
				seen[modelName] = struct{}{}
				out = append(out, modelName)
			}
		}
	}

	return out
}

// collectKnowledgeBases 收集执行中用到的知识库
func collectKnowledgeBases(results []agent.StepResult) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0)

	for _, r := range results {
		if !r.Success {
			continue
		}
		if kb, ok := r.Output["knowledge_base"].(string); ok && kb != "" {
			if _, exists := seen[kb]; !exists {
				seen[kb] = struct{}{}
				out = append(out, kb)
			}
		}
	}

	return out
}

// collectCapabilities 收集执行中用到的能力
func collectCapabilities(results []agent.StepResult) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0)

	for _, r := range results {
		if !r.Success {
			continue
		}
		if capName, ok := r.Output["capability_name"].(string); ok && capName != "" {
			if _, exists := seen[capName]; !exists {
				seen[capName] = struct{}{}
				out = append(out, capName)
			}
		}
	}

	return out
}

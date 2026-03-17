package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/model"
)

// ModelUsageRepository 定义模型用量记录仓储接口
type ModelUsageRepository interface {
	Save(ctx context.Context, record model.UsageRecord) error
}

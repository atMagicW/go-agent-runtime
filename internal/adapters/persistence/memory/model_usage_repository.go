package memory

import (
	"context"
	"sync"

	"github.com/atMagicW/go-agent-runtime/internal/domain/model"
)

// ModelUsageRepository 是模型调用记录的内存实现
type ModelUsageRepository struct {
	mu      sync.Mutex
	records []model.UsageRecord
}

// NewModelUsageRepository 创建内存版模型调用仓储
func NewModelUsageRepository() *ModelUsageRepository {
	return &ModelUsageRepository{
		records: make([]model.UsageRecord, 0),
	}
}

// Save 保存模型调用记录
func (r *ModelUsageRepository) Save(_ context.Context, record model.UsageRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = append(r.records, record)
	return nil
}

package memory

import (
	"context"
	"sync"

	"github.com/atMagicW/go-agent-runtime/internal/domain/audit"
)

// AuditRepository 是审计日志的内存实现
type AuditRepository struct {
	mu      sync.Mutex
	records []audit.Record
}

// NewAuditRepository 创建内存版审计仓储
func NewAuditRepository() *AuditRepository {
	return &AuditRepository{
		records: make([]audit.Record, 0),
	}
}

// Save 保存审计记录
func (r *AuditRepository) Save(_ context.Context, record audit.Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = append(r.records, record)
	return nil
}

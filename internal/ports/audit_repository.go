package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/audit"
)

// AuditRepository 定义审计日志仓储接口
type AuditRepository interface {
	Save(ctx context.Context, record audit.Record) error
}

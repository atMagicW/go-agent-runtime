package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
)

// PromptRepository 定义 Prompt 模板仓储接口
type PromptRepository interface {
	GetByNameAndVersion(ctx context.Context, promptName string, version string) (prompt.Template, error)
	GetLatestByName(ctx context.Context, promptName string) (prompt.Template, error)
	ListByName(ctx context.Context, promptName string) ([]prompt.Template, error)
}

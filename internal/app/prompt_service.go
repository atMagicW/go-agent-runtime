package app

import (
	"context"

	domainprompt "github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// PromptService 提供 Prompt 模板查询服务
type PromptService struct {
	repo ports.PromptRepository
}

// NewPromptService 创建 PromptService
func NewPromptService(repo ports.PromptRepository) *PromptService {
	return &PromptService{
		repo: repo,
	}
}

// GetLatest 获取最新版本 Prompt
func (s *PromptService) GetLatest(ctx context.Context, promptName string) (domainprompt.Template, error) {
	return s.repo.GetLatestByName(ctx, promptName)
}

// GetByVersion 获取指定版本 Prompt
func (s *PromptService) GetByVersion(ctx context.Context, promptName string, version string) (domainprompt.Template, error) {
	return s.repo.GetByNameAndVersion(ctx, promptName, version)
}

// ListVersions 列出所有版本
func (s *PromptService) ListVersions(ctx context.Context, promptName string) ([]domainprompt.Template, error) {
	return s.repo.ListByName(ctx, promptName)
}

package app

import (
	"context"

	pgrag "github.com/atMagicW/go-agent-runtime/internal/adapters/rag/pgvector"
	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// InitKnowledgeBases 根据配置初始化知识库
func InitKnowledgeBases(
	ctx context.Context,
	repo *pgrag.Repository,
	embedding ports.EmbeddingProvider,
	cfgs *cfg.KnowledgeBasesConfig,
) error {
	if cfgs == nil {
		return nil
	}

	for _, item := range cfgs.KnowledgeBases {
		if !item.Enabled {
			continue
		}

		err := repo.EnsureKnowledgeBase(ctx, rag.KnowledgeBase{
			KBID:        item.KBID,
			TenantID:    item.TenantID,
			Name:        item.Name,
			Description: item.Description,
			Enabled:     item.Enabled,
		})
		if err != nil {
			return err
		}

		if item.SeedDemoData {
			if err := pgrag.SeedDemoKnowledgeBase(ctx, repo, embedding, item.KBID); err != nil {
				// demo seed 失败不阻断主流程时，你也可以改成 continue
				return err
			}
		}
	}

	return nil
}

package app

import "github.com/atMagicW/go-agent-runtime/internal/domain/capability"

// SkillService 提供 Skill 查询能力
type SkillService struct {
	registry *SkillRegistry
}

// NewSkillService 创建 SkillService
func NewSkillService(registry *SkillRegistry) *SkillService {
	return &SkillService{
		registry: registry,
	}
}

// ListSkills 列出全部 Skill
func (s *SkillService) ListSkills() []capability.SkillDefinition {
	if s == nil || s.registry == nil {
		return nil
	}
	return s.registry.List()
}

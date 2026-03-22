package app

import (
	"sync"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

// SkillRegistry 管理文件化 Skill 定义
type SkillRegistry struct {
	mu     sync.RWMutex
	skills map[string]capability.SkillDefinition
}

// NewSkillRegistry 创建 SkillRegistry
func NewSkillRegistry(defs []capability.SkillDefinition) *SkillRegistry {
	r := &SkillRegistry{
		skills: make(map[string]capability.SkillDefinition),
	}

	for _, item := range defs {
		r.skills[item.Name] = item
	}

	return r
}

// Get 获取指定 skill
func (r *SkillRegistry) Get(name string) (capability.SkillDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.skills[name]
	return item, ok
}

// List 列出全部 skill
func (r *SkillRegistry) List() []capability.SkillDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]capability.SkillDefinition, 0, len(r.skills))
	for _, item := range r.skills {
		out = append(out, item)
	}
	return out
}

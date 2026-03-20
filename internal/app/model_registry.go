package app

import (
	"strings"

	domainmodel "github.com/atMagicW/go-agent-runtime/internal/domain/model"
	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
)

// ModelRegistry 管理模型配置
type ModelRegistry struct {
	defaultModel  string
	taskTypeToTag map[string]string
	byName        map[string]domainmodel.Profile
}

// NewModelRegistry 从配置构建模型注册表
func NewModelRegistry(c *cfg.ModelsConfig) *ModelRegistry {
	r := &ModelRegistry{
		byName:        make(map[string]domainmodel.Profile),
		taskTypeToTag: map[string]string{},
	}

	if c == nil {
		return r
	}

	r.defaultModel = c.DefaultModel

	for k, v := range c.TaskTypeToTag {
		r.taskTypeToTag[k] = v
	}

	for _, item := range c.Models {
		r.byName[item.Name] = domainmodel.Profile{
			Name:     item.Name,
			Provider: item.Provider,
			Enabled:  item.Enabled,
			Tags:     item.Tags,
		}
	}

	return r
}

// DefaultModel 返回默认模型
func (r *ModelRegistry) DefaultModel() string {
	return r.defaultModel
}

// Get 获取指定模型
func (r *ModelRegistry) Get(name string) (domainmodel.Profile, bool) {
	item, ok := r.byName[name]
	return item, ok
}

// IsEnabled 判断模型是否启用
func (r *ModelRegistry) IsEnabled(name string) bool {
	item, ok := r.byName[name]
	return ok && item.Enabled
}

// ProviderOf 返回模型 provider
func (r *ModelRegistry) ProviderOf(name string) string {
	item, ok := r.byName[name]
	if !ok {
		return ""
	}
	return item.Provider
}

// ResolveByTaskType 根据任务类型选择一个启用模型
func (r *ModelRegistry) ResolveByTaskType(taskType string) (domainmodel.Profile, bool) {
	tag := r.taskTypeToTag[taskType]
	if tag != "" {
		for _, item := range r.byName {
			if !item.Enabled {
				continue
			}
			if contains(item.Tags, tag) {
				return item, true
			}
		}
	}

	if r.defaultModel != "" {
		if item, ok := r.byName[r.defaultModel]; ok && item.Enabled {
			return item, true
		}
	}

	return domainmodel.Profile{}, false
}

// AllEnabledNames 返回所有启用模型名
func (r *ModelRegistry) AllEnabledNames() []string {
	out := make([]string, 0)
	for _, item := range r.byName {
		if item.Enabled {
			out = append(out, item.Name)
		}
	}
	return out
}

func contains(items []string, target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return false
	}

	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

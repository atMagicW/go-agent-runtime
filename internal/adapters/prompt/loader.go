package prompt

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	domainprompt "github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
)

// InMemoryRepository 是第一版 Prompt 模板仓储实现
type InMemoryRepository struct {
	mu        sync.RWMutex
	templates map[string][]domainprompt.Template
}

// NewInMemoryRepository 创建内存版 Prompt 仓储
func NewInMemoryRepository() *InMemoryRepository {
	repo := &InMemoryRepository{
		templates: make(map[string][]domainprompt.Template),
	}

	// 初始化内置模板
	repo.seed()

	return repo
}

// GetByNameAndVersion 获取指定名称和版本的模板
func (r *InMemoryRepository) GetByNameAndVersion(
	_ context.Context,
	promptName string,
	version string,
) (domainprompt.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.templates[promptName]
	for _, item := range items {
		if item.Version == version {
			return item, nil
		}
	}

	return domainprompt.Template{}, errors.New("prompt template not found")
}

// GetLatestByName 获取最新版本模板
func (r *InMemoryRepository) GetLatestByName(
	_ context.Context,
	promptName string,
) (domainprompt.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.templates[promptName]
	if len(items) == 0 {
		return domainprompt.Template{}, errors.New("prompt template not found")
	}

	// 第一版按版本字符串排序，约定 v1 < v2 < v3
	sort.Slice(items, func(i, j int) bool {
		return items[i].Version > items[j].Version
	})

	return items[0], nil
}

// ListByName 列出同名模板所有版本
func (r *InMemoryRepository) ListByName(
	_ context.Context,
	promptName string,
) ([]domainprompt.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.templates[promptName]
	if len(items) == 0 {
		return nil, errors.New("prompt template not found")
	}

	out := make([]domainprompt.Template, len(items))
	copy(out, items)

	sort.Slice(out, func(i, j int) bool {
		return out[i].Version > out[j].Version
	})

	return out, nil
}

// seed 初始化内置 Prompt 模板
func (r *InMemoryRepository) seed() {
	now := time.Now()

	r.templates["intent_classifier"] = []domainprompt.Template{
		{
			PromptName: "intent_classifier",
			Version:    "v1",
			Scene:      "intent",
			Content:    "你是一个意图识别器。请根据用户输入识别意图类型。",
			Variables:  []string{"message"},
			Status:     "active",
			CreatedAt:  now,
		},
	}

	r.templates["final_response"] = []domainprompt.Template{
		{
			PromptName: "final_response",
			Version:    "v1",
			Scene:      "response",
			Content:    "请根据用户请求和前序步骤结果生成最终回答。",
			Variables:  []string{"message", "step_results"},
			Status:     "active",
			CreatedAt:  now,
		},
		{
			PromptName: "final_response",
			Version:    "v2",
			Scene:      "response",
			Content:    "请结合用户请求、知识库证据、工具结果，生成结构化且清晰的最终回答。",
			Variables:  []string{"message", "step_results", "evidences"},
			Status:     "active",
			CreatedAt:  now.Add(time.Minute),
		},
	}
}

package capability

import (
	"fmt"
	"sort"
	"sync"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// Registry 统一管理所有本地能力
type Registry struct {
	mu           sync.RWMutex
	capabilities map[string]ports.Capability
}

// NewRegistry 创建能力注册表
func NewRegistry() *Registry {
	return &Registry{
		capabilities: make(map[string]ports.Capability),
	}
}

// Register 注册一个能力
func (r *Registry) Register(cap ports.Capability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	desc := cap.Descriptor()
	if desc.Name == "" {
		return fmt.Errorf("capability name is empty")
	}

	if _, exists := r.capabilities[desc.Name]; exists {
		return fmt.Errorf("capability already registered: %s", desc.Name)
	}

	r.capabilities[desc.Name] = cap
	return nil
}

// Get 按名称获取能力
func (r *Registry) Get(name string) (ports.Capability, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.capabilities[name]
	return c, ok
}

// ListDescriptors 列出所有能力描述
func (r *Registry) ListDescriptors() []capability.Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]capability.Descriptor, 0, len(r.capabilities))
	for _, cap := range r.capabilities {
		out = append(out, cap.Descriptor())
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out
}

// MustRegister 批量注册，出错直接 panic，适合启动阶段
func (r *Registry) MustRegister(cap ports.Capability) {
	if err := r.Register(cap); err != nil {
		panic(err)
	}
}

package governance

import (
	"sync"
	"time"
)

// BreakerRegistry 用来统一管理不同资源对应的熔断器
type BreakerRegistry struct {
	mu       sync.Mutex
	breakers map[string]*CircuitBreaker
}

// NewBreakerRegistry 创建熔断器注册表
func NewBreakerRegistry() *BreakerRegistry {
	return &BreakerRegistry{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate 获取或创建熔断器
func (r *BreakerRegistry) GetOrCreate(name string, failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	r.mu.Lock()
	defer r.mu.Unlock()

	if b, ok := r.breakers[name]; ok {
		return b
	}

	b := NewCircuitBreaker(name, failureThreshold, resetTimeout)
	r.breakers[name] = b
	return b
}

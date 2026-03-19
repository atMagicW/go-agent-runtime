package governance

import (
	"sync"
	"time"
)

// CircuitState 表示熔断器状态
type CircuitState string

const (
	StateClosed   CircuitState = "closed"
	StateOpen     CircuitState = "open"
	StateHalfOpen CircuitState = "half_open"
)

// CircuitBreaker 是一个轻量级熔断器
type CircuitBreaker struct {
	mu sync.Mutex

	name string

	// 连续失败阈值
	failureThreshold int

	// 熔断打开后，多久进入 half-open
	resetTimeout time.Duration

	state CircuitState

	consecutiveFailures int
	lastFailureTime     time.Time
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(name string, failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	if failureThreshold <= 0 {
		failureThreshold = 3
	}
	if resetTimeout <= 0 {
		resetTimeout = 10 * time.Second
	}

	return &CircuitBreaker{
		name:             name,
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		state:            StateClosed,
	}
}

// Allow 判断当前是否允许请求通过
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// 超过 resetTimeout 后，进入 half-open，允许一次试探请求
		if time.Since(cb.lastFailureTime) >= cb.resetTimeout {
			cb.state = StateHalfOpen
			return true
		}
		return false

	case StateHalfOpen:
		// half-open 只允许单次试探，这里简单允许，由调用后的 success/failure 决定状态
		return true

	default:
		return true
	}
}

// OnSuccess 记录成功
func (cb *CircuitBreaker) OnSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.consecutiveFailures = 0
}

// OnFailure 记录失败
func (cb *CircuitBreaker) OnFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.consecutiveFailures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateHalfOpen:
		cb.state = StateOpen
	case StateClosed:
		if cb.consecutiveFailures >= cb.failureThreshold {
			cb.state = StateOpen
		}
	}
}

// State 获取当前状态
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Name 获取熔断器名称
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

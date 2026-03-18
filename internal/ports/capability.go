package ports

import (
	"context"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

// Capability 定义统一能力接口
type Capability interface {
	// Descriptor 返回能力元信息
	Descriptor() capability.Descriptor

	// Invoke 执行能力
	Invoke(ctx context.Context, input map[string]any) (capability.Result, error)
}

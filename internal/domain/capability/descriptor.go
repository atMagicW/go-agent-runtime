package capability

// Kind 表示能力类型
type Kind string

const (
	KindSkill   Kind = "skill"
	KindTool    Kind = "tool"
	KindMCPTool Kind = "mcp_tool"
)

// Descriptor 描述一个能力的元信息
type Descriptor struct {
	Name        string `json:"name"`
	Kind        Kind   `json:"kind"`
	Description string `json:"description"`

	// 标签，便于后续做路由与筛选
	Tags []string `json:"tags,omitempty"`

	// 版本号，第一版先保留字段
	Version string `json:"version,omitempty"`

	// 是否启用
	Enabled bool `json:"enabled"`
}

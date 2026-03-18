package capability

// Result 表示能力执行结果
type Result struct {
	Name    string         `json:"name"`
	Kind    Kind           `json:"kind"`
	Success bool           `json:"success"`
	Output  map[string]any `json:"output,omitempty"`
	Error   string         `json:"error,omitempty"`
}

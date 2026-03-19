package capability

// MCPToolSpec 描述一个 MCP Tool 的绑定信息
type MCPToolSpec struct {
	Name        string `json:"name"`
	ServerName  string `json:"server_name"`
	RemoteTool  string `json:"remote_tool"`
	Description string `json:"description"`
	Version     string `json:"version,omitempty"`
	Enabled     bool   `json:"enabled"`
}

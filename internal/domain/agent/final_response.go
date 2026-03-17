package agent

// FinalResponse 最终返回结果
type FinalResponse struct {
	Message string `json:"message"`

	Cost float64 `json:"cost,omitempty"`

	Tokens int `json:"tokens,omitempty"`
}

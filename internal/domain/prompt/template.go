package prompt

import "time"

// Template 表示一个 Prompt 模板版本
type Template struct {
	PromptName string `json:"prompt_name"`

	Version string `json:"version"`

	Scene string `json:"scene"`

	Content string `json:"content"`

	// 变量名列表
	Variables []string `json:"variables,omitempty"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

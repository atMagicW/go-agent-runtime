package model

// Profile 表示运行时模型画像
type Profile struct {
	Name     string   `json:"name"`
	Provider string   `json:"provider"`
	Enabled  bool     `json:"enabled"`
	Tags     []string `json:"tags,omitempty"`
}

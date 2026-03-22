package capability

type SkillDefinition struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Entrypoint  string   `json:"entrypoint"`
	Kind        string   `json:"kind"`
	Tags        []string `json:"tags,omitempty"`
	Content     string   `json:"content"`
}

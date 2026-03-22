package skillloader

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

type manifest struct {
	Name        string   `yaml:"name"`
	DisplayName string   `yaml:"display_name"`
	Description string   `yaml:"description"`
	Enabled     bool     `yaml:"enabled"`
	Entrypoint  string   `yaml:"entrypoint"`
	Kind        string   `yaml:"kind"`
	Tags        []string `yaml:"tags"`
}

// FileLoader 从 skills 目录加载 skill 定义
type FileLoader struct {
	baseDir string
}

// NewFileLoader 创建 Skill 文件加载器
func NewFileLoader(baseDir string) *FileLoader {
	return &FileLoader{
		baseDir: baseDir,
	}
}

// Load 加载所有 skill 定义
func (l *FileLoader) Load() ([]capability.SkillDefinition, error) {
	entries, err := os.ReadDir(l.baseDir)
	if err != nil {
		return nil, err
	}

	out := make([]capability.SkillDefinition, 0)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dir := filepath.Join(l.baseDir, entry.Name())

		manifestPath := filepath.Join(dir, "manifest.yaml")
		rawManifest, err := os.ReadFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("read manifest failed: %w", err)
		}

		var m manifest
		if err := yaml.Unmarshal(rawManifest, &m); err != nil {
			return nil, fmt.Errorf("unmarshal manifest failed: %w", err)
		}

		entrypoint := m.Entrypoint
		if entrypoint == "" {
			entrypoint = "SKILL.md"
		}

		rawSkill, err := os.ReadFile(filepath.Join(dir, entrypoint))
		if err != nil {
			return nil, fmt.Errorf("read skill md failed: %w", err)
		}

		out = append(out, capability.SkillDefinition{
			Name:        m.Name,
			DisplayName: m.DisplayName,
			Description: m.Description,
			Enabled:     m.Enabled,
			Entrypoint:  entrypoint,
			Kind:        m.Kind,
			Tags:        m.Tags,
			Content:     string(rawSkill),
		})
	}

	return out, nil
}

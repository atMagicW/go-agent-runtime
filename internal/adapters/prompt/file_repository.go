package prompt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	domainprompt "github.com/atMagicW/go-agent-runtime/internal/domain/prompt"
)

// FileRepository 是基于 prompts 目录的模板仓储实现
type FileRepository struct {
	mu        sync.RWMutex
	baseDir   string
	templates map[string][]domainprompt.Template
}

// NewFileRepository 创建文件版 Prompt 仓储
func NewFileRepository(baseDir string) (*FileRepository, error) {
	r := &FileRepository{
		baseDir:   baseDir,
		templates: make(map[string][]domainprompt.Template),
	}

	if err := r.Load(); err != nil {
		return nil, err
	}

	return r, nil
}

// Load 扫描 prompts 目录并加载模板
func (r *FileRepository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make(map[string][]domainprompt.Template)

	err := filepath.Walk(r.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".tmpl") {
			return nil
		}

		tpl, parseErr := parseTemplateFile(path, info)
		if parseErr != nil {
			return parseErr
		}

		items[tpl.PromptName] = append(items[tpl.PromptName], tpl)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk prompt dir failed: %w", err)
	}

	// 每个 prompt_name 下按版本倒序
	for name, versions := range items {
		sort.Slice(versions, func(i, j int) bool {
			return compareVersionDesc(versions[i].Version, versions[j].Version)
		})
		items[name] = versions
	}

	r.templates = items
	return nil
}

// GetByNameAndVersion 获取指定名称和版本模板
func (r *FileRepository) GetByNameAndVersion(
	_ context.Context,
	promptName string,
	version string,
) (domainprompt.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.templates[promptName]
	for _, item := range items {
		if item.Version == version {
			return item, nil
		}
	}

	return domainprompt.Template{}, fmt.Errorf("prompt template not found: %s@%s", promptName, version)
}

// GetLatestByName 获取最新版本模板
func (r *FileRepository) GetLatestByName(
	_ context.Context,
	promptName string,
) (domainprompt.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.templates[promptName]
	if len(items) == 0 {
		return domainprompt.Template{}, fmt.Errorf("prompt template not found: %s", promptName)
	}

	return items[0], nil
}

// ListByName 列出某个 prompt 的全部版本
func (r *FileRepository) ListByName(
	_ context.Context,
	promptName string,
) ([]domainprompt.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.templates[promptName]
	if len(items) == 0 {
		return nil, fmt.Errorf("prompt template not found: %s", promptName)
	}

	out := make([]domainprompt.Template, len(items))
	copy(out, items)
	return out, nil
}

func parseTemplateFile(path string, info os.FileInfo) (domainprompt.Template, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return domainprompt.Template{}, fmt.Errorf("read prompt file failed: %w", err)
	}

	fileName := strings.TrimSuffix(info.Name(), ".tmpl")
	promptName, version, err := splitPromptFileName(fileName)
	if err != nil {
		return domainprompt.Template{}, err
	}

	scene := filepath.Base(filepath.Dir(path))

	return domainprompt.Template{
		PromptName: promptName,
		Version:    version,
		Scene:      scene,
		Content:    string(content),
		Variables:  nil,
		Status:     "active",
		CreatedAt:  info.ModTime(),
	}, nil
}

// splitPromptFileName 解析类似 final_response_v2
func splitPromptFileName(fileName string) (string, string, error) {
	idx := strings.LastIndex(fileName, "_v")
	if idx <= 0 || idx >= len(fileName)-2 {
		return "", "", fmt.Errorf("invalid prompt file name: %s", fileName)
	}

	promptName := fileName[:idx]
	version := fileName[idx+1:] // 保留 v2 这种格式

	return promptName, version, nil
}

// compareVersionDesc 简单比较 v10 > v2
func compareVersionDesc(a, b string) bool {
	ai := parseVersionNum(a)
	bi := parseVersionNum(b)
	if ai == bi {
		return a > b
	}
	return ai > bi
}

func parseVersionNum(v string) int {
	v = strings.TrimSpace(strings.TrimPrefix(v, "v"))
	n := 0
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			break
		}
		n = n*10 + int(ch-'0')
	}
	return n
}

// ReloadForTestOrHotUpdate 提供后续热更新入口
func (r *FileRepository) ReloadForTestOrHotUpdate() error {
	return r.Load()
}

// Snapshot 返回当前加载结果，便于调试
func (r *FileRepository) Snapshot() map[string][]domainprompt.Template {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[string][]domainprompt.Template, len(r.templates))
	for k, v := range r.templates {
		cp := make([]domainprompt.Template, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

// Ensure import time used
var _ = time.Now

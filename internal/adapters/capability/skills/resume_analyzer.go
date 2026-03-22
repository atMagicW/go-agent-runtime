package skills

import (
	"context"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

// ResumeAnalyzerSkill 是一个本地简历分析 Skill
type ResumeAnalyzerSkill struct {
	def capability.SkillDefinition
}

// NewResumeAnalyzerSkill 创建简历分析 Skill
func NewResumeAnalyzerSkill(def capability.SkillDefinition) *ResumeAnalyzerSkill {
	return &ResumeAnalyzerSkill{
		def: def,
	}
}

// Descriptor 返回 Skill 元信息
func (s *ResumeAnalyzerSkill) Descriptor() capability.Descriptor {
	description := "分析简历文本，提取候选人的核心优势与建议"
	tags := []string{"resume", "analysis", "skill"}
	enabled := true

	if s.def.Name != "" {
		if s.def.Description != "" {
			description = s.def.Description
		}
		if len(s.def.Tags) > 0 {
			tags = s.def.Tags
		}
		enabled = s.def.Enabled
	}

	return capability.Descriptor{
		Name:        "resume_analyzer",
		Kind:        capability.KindSkill,
		Description: description,
		Tags:        tags,
		Version:     "v1",
		Enabled:     enabled,
	}
}

// Invoke 执行 Skill
func (s *ResumeAnalyzerSkill) Invoke(_ context.Context, input map[string]any) (capability.Result, error) {
	text, _ := input["message"].(string)

	advantages := make([]string, 0)
	suggestions := make([]string, 0)

	lower := strings.ToLower(text)

	if strings.Contains(lower, "golang") || strings.Contains(lower, "go") {
		advantages = append(advantages, "具备 Go 相关经验，可作为后端/基础设施方向亮点")
	}
	if strings.Contains(lower, "agent") {
		advantages = append(advantages, "具备 Agent 系统设计意识，适合 AI 应用工程岗位")
	}
	if strings.Contains(lower, "rag") {
		advantages = append(advantages, "具备 RAG 相关认知，适合知识系统与检索增强场景")
	}

	if len(advantages) == 0 {
		advantages = append(advantages, "建议突出技术栈、项目结果和业务价值")
	}

	suggestions = append(suggestions, "建议在简历中突出项目中的架构设计、模块职责与工程化能力")
	suggestions = append(suggestions, "建议量化结果，例如性能提升、时延降低、成本优化等")

	return capability.Result{
		Name:    "resume_analyzer",
		Kind:    capability.KindSkill,
		Success: true,
		Output: map[string]any{
			"capability_name": "resume_analyzer",
			"kind":            "skill",
			"skill_prompt":    s.def.Content,
			"advantages":      advantages,
			"suggestions":     suggestions,
			"result":          "已完成简历分析",
		},
	}, nil
}

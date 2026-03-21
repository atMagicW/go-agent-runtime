package response

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// TemplateResponseComposer 是基于 Prompt 模板的最终回答生成器
type TemplateResponseComposer struct {
	promptRepo ports.PromptRepository
}

// NewTemplateResponseComposer 创建 ResponseComposer
func NewTemplateResponseComposer(promptRepo ports.PromptRepository) *TemplateResponseComposer {
	return &TemplateResponseComposer{
		promptRepo: promptRepo,
	}
}

// Compose 生成最终回答
func (c *TemplateResponseComposer) Compose(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ComposeRequest,
) (ports.ComposeResponse, error) {
	renderedPrompt, err := c.BuildPrompt(ctx, runtimeCtx, req)
	if err != nil {
		return ports.ComposeResponse{}, err
	}

	text := c.extractSummary(req.StepResults)
	if text == "" {
		text = renderedPrompt
	}

	return ports.ComposeResponse{
		Text:   text,
		Tokens: len([]rune(text)) / 4,
		Cost:   0.0005,
		Model:  "template-composer",
	}, nil
}

// buildStepResultsText 将步骤结果转成可读文本
func (c *TemplateResponseComposer) buildStepResultsText(results []agent.StepResult) string {
	if len(results) == 0 {
		return "无"
	}

	var sb strings.Builder
	for i, result := range results {
		sb.WriteString(fmt.Sprintf("%d. step_id=%s success=%v\n", i+1, result.StepID, result.Success))

		if result.Error != "" {
			sb.WriteString("   error=")
			sb.WriteString(result.Error)
			sb.WriteString("\n")
		}

		if len(result.Output) > 0 {
			raw, _ := json.Marshal(result.Output)
			sb.WriteString("   output=")
			sb.WriteString(string(raw))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// buildEvidencesText 提取检索证据
func (c *TemplateResponseComposer) buildEvidencesText(results []agent.StepResult) string {
	lines := make([]string, 0)

	for _, result := range results {
		if !result.Success {
			continue
		}

		raw, ok := result.Output["evidences"]
		if !ok {
			continue
		}

		items, ok := raw.([]map[string]any)
		if ok {
			for _, item := range items {
				content, _ := item["content"].(string)
				kb, _ := item["kb"].(string)
				if content != "" {
					lines = append(lines, fmt.Sprintf("[kb=%s] %s", kb, content))
				}
			}
			continue
		}

		// 兼容经过 interface{} 传递后的情况
		if list, ok := raw.([]any); ok {
			for _, one := range list {
				if item, ok := one.(map[string]any); ok {
					content, _ := item["content"].(string)
					kb, _ := item["kb"].(string)
					if content != "" {
						lines = append(lines, fmt.Sprintf("[kb=%s] %s", kb, content))
					}
				}
			}
		}
	}

	if len(lines) == 0 {
		return "无"
	}

	return strings.Join(lines, "\n")
}

// buildCapabilityResultsText 提取能力结果
func (c *TemplateResponseComposer) buildCapabilityResultsText(results []agent.StepResult) string {
	lines := make([]string, 0)

	for _, result := range results {
		if !result.Success {
			continue
		}

		name, _ := result.Output["capability_name"].(string)
		kind, _ := result.Output["kind"].(string)

		if name == "" {
			continue
		}

		if summary, ok := result.Output["result"].(string); ok && summary != "" {
			lines = append(lines, fmt.Sprintf("[%s/%s] %s", kind, name, summary))
			continue
		}

		raw, _ := json.Marshal(result.Output)
		lines = append(lines, fmt.Sprintf("[%s/%s] %s", kind, name, string(raw)))
	}

	if len(lines) == 0 {
		return "无"
	}

	return strings.Join(lines, "\n")
}

// buildFallbackPrompt 模板渲染失败时使用的降级文本
func (c *TemplateResponseComposer) buildFallbackPrompt(data TemplateData) string {
	var sb strings.Builder

	sb.WriteString("用户请求：\n")
	sb.WriteString(data.Message)
	sb.WriteString("\n\n")

	sb.WriteString("识别意图：\n")
	sb.WriteString(data.Intent)
	sb.WriteString("\n\n")

	sb.WriteString("步骤结果：\n")
	sb.WriteString(data.StepResultsText)
	sb.WriteString("\n\n")

	if data.EvidencesText != "无" {
		sb.WriteString("知识证据：\n")
		sb.WriteString(data.EvidencesText)
		sb.WriteString("\n\n")
	}

	if data.CapabilityResultsText != "无" {
		sb.WriteString("能力结果：\n")
		sb.WriteString(data.CapabilityResultsText)
		sb.WriteString("\n\n")
	}

	sb.WriteString("请基于以上信息生成最终回答。")
	return sb.String()
}

func (c *TemplateResponseComposer) BuildPrompt(
	ctx context.Context,
	runtimeCtx agent.RuntimeContext,
	req ports.ComposeRequest,
) (string, error) {
	var promptContent string

	if req.PromptVer != "" {
		tpl, err := c.promptRepo.GetByNameAndVersion(ctx, req.PromptName, req.PromptVer)
		if err != nil {
			return "", err
		}
		promptContent = tpl.Content
	} else {
		tpl, err := c.promptRepo.GetLatestByName(ctx, req.PromptName)
		if err != nil {
			return "", err
		}
		promptContent = tpl.Content
	}

	data := TemplateData{
		Message:               req.Message,
		Intent:                string(runtimeCtx.Intent.IntentType),
		StepResultsText:       c.buildStepResultsText(req.StepResults),
		EvidencesText:         c.buildEvidencesText(req.StepResults),
		CapabilityResultsText: c.buildCapabilityResultsText(req.StepResults),
	}

	renderedPrompt, err := RenderTemplate(promptContent, data)
	if err != nil {
		return c.buildFallbackPrompt(data), nil
	}

	return renderedPrompt, nil
}

// extractSummary 优先抽取适合作为最终回答的文本
func (c *TemplateResponseComposer) extractSummary(results []agent.StepResult) string {
	for i := len(results) - 1; i >= 0; i-- {
		r := results[i]
		if !r.Success {
			continue
		}

		if text, ok := r.Output["text"].(string); ok && text != "" {
			return text
		}

		if text, ok := r.Output["result"].(string); ok && text != "" {
			return text
		}
	}

	return ""
}

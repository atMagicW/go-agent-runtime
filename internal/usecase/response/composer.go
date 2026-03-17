package response

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/agent"
	"github.com/atMagicW/go-agent-runtime/internal/ports"
)

// TemplateResponseComposer 是基于 Prompt 模板的第一版最终回答生成器
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
	var (
		promptContent string
	)

	// 优先按指定版本获取，否则取最新版本
	if req.PromptVer != "" {
		tpl, getErr := c.promptRepo.GetByNameAndVersion(ctx, req.PromptName, req.PromptVer)
		if getErr != nil {
			return ports.ComposeResponse{}, getErr
		}
		promptContent = tpl.Content
	} else {
		tpl, getErr := c.promptRepo.GetLatestByName(ctx, req.PromptName)
		if getErr != nil {
			return ports.ComposeResponse{}, getErr
		}
		promptContent = tpl.Content
	}

	// 第一版先不做真正模板渲染引擎，
	// 而是把模板文本 + 用户消息 + step 结果拼接成可读结果。
	text := c.buildFinalText(promptContent, runtimeCtx, req)

	return ports.ComposeResponse{
		Text:   text,
		Tokens: len([]rune(text)) / 4,
		Cost:   0.0005,
		Model:  "template-composer",
	}, nil
}

// buildFinalText 构造最终文本
func (c *TemplateResponseComposer) buildFinalText(
	promptContent string,
	runtimeCtx agent.RuntimeContext,
	req ports.ComposeRequest,
) string {
	var sb strings.Builder

	sb.WriteString("【系统模板】\n")
	sb.WriteString(promptContent)
	sb.WriteString("\n\n")

	sb.WriteString("【用户请求】\n")
	sb.WriteString(req.Message)
	sb.WriteString("\n\n")

	if runtimeCtx.Intent.IntentType != "" {
		sb.WriteString("【识别意图】\n")
		sb.WriteString(string(runtimeCtx.Intent.IntentType))
		sb.WriteString("\n\n")
	}

	sb.WriteString("【执行结果汇总】\n")

	for i, result := range req.StepResults {
		sb.WriteString(fmt.Sprintf("%d. 步骤ID=%s\n", i+1, result.StepID))
		sb.WriteString(fmt.Sprintf("   成功=%v\n", result.Success))

		if result.Error != "" {
			sb.WriteString(fmt.Sprintf("   错误=%s\n", result.Error))
		}

		if len(result.Output) > 0 {
			raw, _ := json.Marshal(result.Output)
			sb.WriteString(fmt.Sprintf("   输出=%s\n", string(raw)))
		}

		sb.WriteString("\n")
	}

	// 额外尝试从步骤结果里提取更“像回答”的内容
	summary := c.extractSummary(req.StepResults)
	if summary != "" {
		sb.WriteString("【最终回答】\n")
		sb.WriteString(summary)
	} else {
		sb.WriteString("【最终回答】\n")
		sb.WriteString("已根据执行步骤完成处理。你可以继续追问更细的内容。")
	}

	return sb.String()
}

// extractSummary 优先抽取更适合作为最终回答的文本
func (c *TemplateResponseComposer) extractSummary(results []agent.StepResult) string {
	// 从后往前找，优先使用最后一个成功步骤的 text / result
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

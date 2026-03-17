package response

import (
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
func (c *TemplateResponseComposer) buildFinalText(
	promptContent string,
	runtimeCtx agent.RuntimeContext,
	req ports.ComposeRequest,
) string {
	summary := c.extractSummary(req.StepResults)

	// 如果已经有比较合适的最终文本，就优先直接返回更自然的用户结果
	if summary != "" {
		return summary
	}

	var sb strings.Builder

	sb.WriteString("我已经根据你的请求完成处理。")

	if runtimeCtx.Intent.IntentType != "" {
		sb.WriteString("本次任务识别到的意图类型为：")
		sb.WriteString(string(runtimeCtx.Intent.IntentType))
		sb.WriteString("。")
	}

	if len(req.StepResults) > 0 {
		sb.WriteString("共执行 ")
		sb.WriteString(fmt.Sprintf("%d", len(req.StepResults)))
		sb.WriteString(" 个步骤。")
	}

	// 如果需要调试信息，可以在后续通过开关再展开
	_ = promptContent

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

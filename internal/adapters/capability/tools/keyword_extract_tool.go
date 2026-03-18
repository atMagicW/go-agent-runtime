package tools

import (
	"context"
	"sort"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

// KeywordExtractTool 是一个本地关键词提取 Tool
type KeywordExtractTool struct {
}

// NewKeywordExtractTool 创建关键词提取 Tool
func NewKeywordExtractTool() *KeywordExtractTool {
	return &KeywordExtractTool{}
}

// Descriptor 返回 Tool 元信息
func (t *KeywordExtractTool) Descriptor() capability.Descriptor {
	return capability.Descriptor{
		Name:        "keyword_extract_tool",
		Kind:        capability.KindTool,
		Description: "从输入文本中提取高频关键词",
		Tags:        []string{"tool", "keyword", "text"},
		Version:     "v1",
		Enabled:     true,
	}
}

// Invoke 执行 Tool
func (t *KeywordExtractTool) Invoke(_ context.Context, input map[string]any) (capability.Result, error) {
	text, _ := input["message"].(string)

	words := strings.Fields(strings.ToLower(text))
	freq := make(map[string]int)

	stopWords := map[string]struct{}{
		"the": {}, "a": {}, "an": {}, "and": {}, "or": {},
		"to": {}, "of": {}, "in": {}, "for": {}, "on": {},
		"我": {}, "你": {}, "他": {}, "她": {}, "它": {},
		"的": {}, "了": {}, "和": {}, "是": {}, "在": {},
		"一个": {}, "这个": {}, "那个": {},
	}

	for _, w := range words {
		w = strings.TrimSpace(w)
		if len([]rune(w)) <= 1 {
			continue
		}
		if _, ok := stopWords[w]; ok {
			continue
		}
		freq[w]++
	}

	type pair struct {
		Key   string
		Count int
	}

	items := make([]pair, 0, len(freq))
	for k, v := range freq {
		items = append(items, pair{Key: k, Count: v})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Key < items[j].Key
		}
		return items[i].Count > items[j].Count
	})

	keywords := make([]string, 0, 5)
	for i, item := range items {
		if i >= 5 {
			break
		}
		keywords = append(keywords, item.Key)
	}

	return capability.Result{
		Name:    "keyword_extract_tool",
		Kind:    capability.KindTool,
		Success: true,
		Output: map[string]any{
			"capability_name": "keyword_extract_tool",
			"kind":            "tool",
			"keywords":        keywords,
			"result":          "已完成关键词提取",
		},
	}, nil
}

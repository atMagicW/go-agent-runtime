package rerank

import (
	"context"
	"sort"
	"strings"

	"github.com/atMagicW/go-agent-runtime/internal/domain/rag"
)

// SimpleReranker 是一个轻量级本地重排器
type SimpleReranker struct {
}

// NewSimpleReranker 创建本地重排器
func NewSimpleReranker() *SimpleReranker {
	return &SimpleReranker{}
}

// Rerank 根据 query 和候选证据进行重排
func (r *SimpleReranker) Rerank(_ context.Context, query string, items []rag.Evidence, topK int) ([]rag.Evidence, error) {
	if len(items) == 0 {
		return items, nil
	}
	if topK <= 0 || topK > len(items) {
		topK = len(items)
	}

	queryTerms := tokenize(query)

	type scored struct {
		item  rag.Evidence
		score float64
	}

	scoredItems := make([]scored, 0, len(items))

	for _, item := range items {
		contentTerms := tokenize(item.Content)

		overlap := calcOverlapScore(queryTerms, contentTerms)

		// 第一版简单加权：
		// 最终分 = 原始召回分 * 0.7 + 词项重叠分 * 0.3
		finalScore := item.Score*0.7 + overlap*0.3
		item.Score = finalScore

		scoredItems = append(scoredItems, scored{
			item:  item,
			score: finalScore,
		})
	}

	sort.Slice(scoredItems, func(i, j int) bool {
		if scoredItems[i].score == scoredItems[j].score {
			return scoredItems[i].item.ChunkID < scoredItems[j].item.ChunkID
		}
		return scoredItems[i].score > scoredItems[j].score
	})

	out := make([]rag.Evidence, 0, topK)
	for i := 0; i < topK; i++ {
		out = append(out, scoredItems[i].item)
	}

	return out, nil
}

func tokenize(text string) []string {
	raw := strings.Fields(strings.ToLower(text))
	out := make([]string, 0, len(raw))

	for _, item := range raw {
		item = strings.TrimSpace(item)
		item = strings.Trim(item, ".,;:!?()[]{}\"'")
		if item == "" {
			continue
		}
		out = append(out, item)
	}

	return out
}

func calcOverlapScore(queryTerms, contentTerms []string) float64 {
	if len(queryTerms) == 0 || len(contentTerms) == 0 {
		return 0
	}

	contentSet := make(map[string]struct{}, len(contentTerms))
	for _, t := range contentTerms {
		contentSet[t] = struct{}{}
	}

	hit := 0
	for _, t := range queryTerms {
		if _, ok := contentSet[t]; ok {
			hit++
		}
	}

	return float64(hit) / float64(len(queryTerms))
}

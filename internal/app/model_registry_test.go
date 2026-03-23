package app

import (
	"testing"

	"github.com/atMagicW/go-agent-runtime/internal/domain/model"
	cfg "github.com/atMagicW/go-agent-runtime/internal/pkg/config"
)

func TestModelRegistry_ResolveByTaskType(t *testing.T) {
	reg := NewModelRegistry(&cfg.ModelsConfig{
		DefaultModel: "gpt-4.1-mini",
		TaskTypeToTag: map[string]string{
			"analysis": "analysis",
			"intent":   "intent",
		},
		Models: []cfg.ModelConfig{
			{
				Name:     "gpt-4.1",
				Provider: string(model.ProviderOpenAI),
				Enabled:  true,
				Tags:     []string{"analysis"},
			},
			{
				Name:     "gpt-4.1-mini",
				Provider: string(model.ProviderOpenAI),
				Enabled:  true,
				Tags:     []string{"intent", "chat"},
			},
		},
	})

	item, ok := reg.ResolveByTaskType("analysis")
	if !ok {
		t.Fatal("ResolveByTaskType() expected ok=true")
	}

	if item.Name != "gpt-4.1" {
		t.Fatalf("ResolveByTaskType() got %s, want gpt-4.1", item.Name)
	}
}

func TestModelRegistry_DefaultModel(t *testing.T) {
	reg := NewModelRegistry(&cfg.ModelsConfig{
		DefaultModel: "gpt-4.1-mini",
		Models: []cfg.ModelConfig{
			{
				Name:     "gpt-4.1-mini",
				Provider: string(model.ProviderOpenAI),
				Enabled:  true,
				Tags:     []string{"chat"},
			},
		},
	})

	if reg.DefaultModel() != "gpt-4.1-mini" {
		t.Fatalf("DefaultModel() got %s", reg.DefaultModel())
	}

	if !reg.IsEnabled("gpt-4.1-mini") {
		t.Fatal("expected model enabled")
	}
}

package capability

import (
	"context"
	"testing"

	domaincap "github.com/atMagicW/go-agent-runtime/internal/domain/capability"
)

type mockCapability struct {
	name string
}

func (m *mockCapability) Descriptor() domaincap.Descriptor {
	return domaincap.Descriptor{
		Name:    m.name,
		Kind:    domaincap.KindTool,
		Enabled: true,
	}
}

func (m *mockCapability) Invoke(ctx context.Context, input map[string]any) (domaincap.Result, error) {
	_ = ctx
	_ = input
	return domaincap.Result{
		Name:    m.name,
		Kind:    domaincap.KindTool,
		Success: true,
		Output:  map[string]any{"ok": true},
	}, nil
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()

	err := r.Register(&mockCapability{name: "mock_tool"})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	got, ok := r.Get("mock_tool")
	if !ok || got == nil {
		t.Fatal("Get() failed")
	}
}

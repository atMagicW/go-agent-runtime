package prompt

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileRepository_LoadAndQuery(t *testing.T) {
	dir := t.TempDir()

	responseDir := filepath.Join(dir, "response")
	if err := os.MkdirAll(responseDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	err := os.WriteFile(filepath.Join(responseDir, "final_response_v1.tmpl"), []byte("v1 content"), 0o644)
	if err != nil {
		t.Fatalf("write v1 failed: %v", err)
	}

	err = os.WriteFile(filepath.Join(responseDir, "final_response_v2.tmpl"), []byte("v2 content"), 0o644)
	if err != nil {
		t.Fatalf("write v2 failed: %v", err)
	}

	repo, err := NewFileRepository(dir)
	if err != nil {
		t.Fatalf("NewFileRepository failed: %v", err)
	}

	latest, err := repo.GetLatestByName(context.Background(), "final_response")
	if err != nil {
		t.Fatalf("GetLatestByName failed: %v", err)
	}

	if latest.Version != "v2" {
		t.Fatalf("latest version = %s, want v2", latest.Version)
	}

	items, err := repo.ListByName(context.Background(), "final_response")
	if err != nil {
		t.Fatalf("ListByName failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
}

package textsplitter

import "testing"

func TestSplit(t *testing.T) {
	s := NewSplitter(10, 2)

	chunks := s.Split("abcdefghijklmnopqrstuvwxyz")
	if len(chunks) < 2 {
		t.Fatalf("Split() chunk count = %d, want >= 2", len(chunks))
	}

	if chunks[0] == "" {
		t.Fatal("first chunk is empty")
	}
}

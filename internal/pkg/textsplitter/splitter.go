package textsplitter

// Splitter 是一个简单文本切块器
type Splitter struct {
	ChunkSize int
	Overlap   int
}

// NewSplitter 创建切块器
func NewSplitter(chunkSize, overlap int) *Splitter {
	if chunkSize <= 0 {
		chunkSize = 300
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 4
	}

	return &Splitter{
		ChunkSize: chunkSize,
		Overlap:   overlap,
	}
}

// Split 将文本切成多个 chunk
func (s *Splitter) Split(text string) []string {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	if len(runes) <= s.ChunkSize {
		return []string{text}
	}

	out := make([]string, 0)
	step := s.ChunkSize - s.Overlap
	if step <= 0 {
		step = s.ChunkSize
	}

	for start := 0; start < len(runes); start += step {
		end := start + s.ChunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunk := string(runes[start:end])
		out = append(out, chunk)

		if end == len(runes) {
			break
		}
	}

	return out
}

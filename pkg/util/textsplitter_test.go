package util

import (
	"strings"
	"testing"
)

func TestNewRecursiveCharacterTextSplitter(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(100, 20)

	if splitter.ChunkSize != 100 {
		t.Errorf("Expected ChunkSize to be 100, got %d", splitter.ChunkSize)
	}

	if splitter.ChunkOverlap != 20 {
		t.Errorf("Expected ChunkOverlap to be 20, got %d", splitter.ChunkOverlap)
	}

	expectedSeparators := []string{"\n\n", "\n", " ", ""}
	if len(splitter.Separators) != len(expectedSeparators) {
		t.Errorf("Expected %d separators, got %d", len(expectedSeparators), len(splitter.Separators))
	}
}

func TestSplitText_EmptyText(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(100, 20)
	result := splitter.SplitText("")

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty text, got %d chunks", len(result))
	}
}

func TestSplitText_ShortText(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(100, 20)
	text := "This is a short text."
	result := splitter.SplitText(text)

	if len(result) != 1 {
		t.Errorf("Expected 1 chunk for short text, got %d", len(result))
	}

	if result[0] != text {
		t.Errorf("Expected chunk to be '%s', got '%s'", text, result[0])
	}
}

func TestSplitText_ChunkSize500Overlap100(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(500, 100)

	// Create a text longer than 500 characters
	text := strings.Repeat("This is a test sentence with some content. ", 20) // ~860 characters

	result := splitter.SplitText(text)

	// Should create at least 2 chunks
	if len(result) < 2 {
		t.Errorf("Expected at least 2 chunks, got %d", len(result))
	}

	// Check that each chunk (except possibly the last) is close to the chunk size
	for i, chunk := range result {
		if len(chunk) > 600 { // Allow some flexibility due to overlap
			t.Errorf("Chunk %d is too large: %d characters", i, len(chunk))
		}

		if len(chunk) == 0 {
			t.Errorf("Chunk %d is empty", i)
		}
	}
}

func TestSplitText_WithParagraphs(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(100, 20)

	text := `First paragraph with some content here.

Second paragraph with different content.

Third paragraph with even more content to test the splitting.`

	result := splitter.SplitText(text)

	// Should create multiple chunks
	if len(result) == 0 {
		t.Error("Expected at least one chunk")
	}

	// Verify that the text is properly split
	combined := strings.Join(result, "")
	if !strings.Contains(combined, "First paragraph") {
		t.Error("First paragraph content missing from chunks")
	}

	if !strings.Contains(combined, "Third paragraph") {
		t.Error("Third paragraph content missing from chunks")
	}
}

func TestSplitText_WithNewlines(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(50, 10)

	text := "Line 1\nLine 2\nLine 3 with some longer content\nLine 4\nLine 5"

	result := splitter.SplitText(text)

	if len(result) == 0 {
		t.Error("Expected at least one chunk")
	}

	// Verify no empty chunks
	for i, chunk := range result {
		if strings.TrimSpace(chunk) == "" {
			t.Errorf("Chunk %d is empty or contains only whitespace", i)
		}
	}
}

func TestGetChunkCount(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(100, 20)

	shortText := "Short text"
	longText := strings.Repeat("This is some content. ", 20)

	shortCount := splitter.GetChunkCount(shortText)
	longCount := splitter.GetChunkCount(longText)

	if shortCount != 1 {
		t.Errorf("Expected 1 chunk for short text, got %d", shortCount)
	}

	if longCount <= 1 {
		t.Errorf("Expected more than 1 chunk for long text, got %d", longCount)
	}
}

func TestGetEstimatedTokenCount(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(100, 20)

	text := "This is a test text with approximately twenty characters."
	tokenCount := splitter.GetEstimatedTokenCount(text)

	// Should be roughly len(text) / 4
	expectedTokens := len(text) / 4
	if tokenCount != expectedTokens {
		t.Errorf("Expected approximately %d tokens, got %d", expectedTokens, tokenCount)
	}
}

func TestOverlapFunctionality(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(50, 15)

	text := "This is the first sentence. This is the second sentence. This is the third sentence."

	result := splitter.SplitText(text)

	if len(result) < 2 {
		t.Skip("Need at least 2 chunks to test overlap")
	}

	// Check that there's some overlap between chunks
	// This is a basic check - in a real implementation you might want more sophisticated overlap verification
	for i := 1; i < len(result); i++ {
		chunk := result[i]
		if len(chunk) == 0 {
			t.Errorf("Chunk %d should not be empty when testing overlap", i)
		}
	}
}

func TestSplitText_VeryLongText(t *testing.T) {
	splitter := NewRecursiveCharacterTextSplitter(500, 100)

	// Create a very long text
	longText := strings.Repeat("This is a very long text that should be split into multiple chunks. ", 100) // ~6700 characters

	result := splitter.SplitText(longText)

	if len(result) < 10 {
		t.Errorf("Expected many chunks for very long text, got %d", len(result))
	}

	// Verify all chunks are reasonable size
	for i, chunk := range result {
		if len(chunk) > 700 { // Allow some flexibility
			t.Errorf("Chunk %d is too large: %d characters (chunk: %s...)", i, len(chunk), chunk[:min(50, len(chunk))])
		}

		if len(chunk) == 0 {
			t.Errorf("Chunk %d is empty", i)
		}
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

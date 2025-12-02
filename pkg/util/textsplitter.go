package util

import (
	"strings"
	"unicode"
)

// TextSplitter defines the interface for text splitting functionality
type TextSplitter interface {
	SplitText(text string) []string
}

// RecursiveCharacterTextSplitter implements a recursive character-based text splitter
type RecursiveCharacterTextSplitter struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string
}

// NewRecursiveCharacterTextSplitter creates a new recursive character text splitter
func NewRecursiveCharacterTextSplitter(chunkSize, chunkOverlap int) *RecursiveCharacterTextSplitter {
	return &RecursiveCharacterTextSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		Separators:   []string{"\n\n", "\n", " ", ""},
	}
}

// SplitText splits the input text into chunks using recursive character splitting
func (r *RecursiveCharacterTextSplitter) SplitText(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	// If text is smaller than chunk size, return as is
	if len(text) <= r.ChunkSize {
		return []string{text}
	}

	return r.splitTextRecursive(text, r.Separators)
}

// splitTextRecursive performs the recursive splitting logic
func (r *RecursiveCharacterTextSplitter) splitTextRecursive(text string, separators []string) []string {
	var finalChunks []string
	separator := ""

	// Find the best separator to use
	for _, sep := range separators {
		if sep == "" {
			separator = sep
			break
		}
		if strings.Contains(text, sep) {
			separator = sep
			break
		}
	}

	var splits []string
	if separator != "" {
		splits = strings.Split(text, separator)
	} else {
		splits = []string{text}
	}

	// Now go through the splits and merge them into good chunks
	goodSplits := []string{}
	for _, split := range splits {
		if len(split) > r.ChunkSize {
			// If we have good splits so far, add them to final chunks
			if len(goodSplits) > 0 {
				mergedChunks := r.mergeSplits(goodSplits, separator)
				finalChunks = append(finalChunks, mergedChunks...)
				goodSplits = []string{}
			}
			// Recursively split the oversized chunk
			finalChunks = append(finalChunks, r.splitTextRecursive(split, separators[1:])...)
		} else {
			goodSplits = append(goodSplits, split)
		}
	}

	// Process remaining good splits
	if len(goodSplits) > 0 {
		mergedChunks := r.mergeSplits(goodSplits, separator)
		finalChunks = append(finalChunks, mergedChunks...)
	}

	return r.addOverlap(finalChunks)
}

// mergeSplits merges the splits while respecting chunk size limits
func (r *RecursiveCharacterTextSplitter) mergeSplits(splits []string, separator string) []string {
	var docs []string
	var currentDoc []string
	total := 0

	for _, split := range splits {
		splitLen := len(split)
		separatorLen := 0
		if len(currentDoc) > 0 {
			separatorLen = len(separator)
		}

		if total+splitLen+separatorLen > r.ChunkSize && len(currentDoc) > 0 {
			// Add current document
			docs = append(docs, strings.Join(currentDoc, separator))
			// Start new document with overlap consideration
			currentDoc = []string{}
			total = 0
		}

		currentDoc = append(currentDoc, split)
		total += splitLen + separatorLen
	}

	// Add the last document if it exists
	if len(currentDoc) > 0 {
		docs = append(docs, strings.Join(currentDoc, separator))
	}

	return docs
}

// addOverlap adds overlap between chunks
func (r *RecursiveCharacterTextSplitter) addOverlap(chunks []string) []string {
	if len(chunks) <= 1 || r.ChunkOverlap <= 0 {
		return chunks
	}

	result := make([]string, 0, len(chunks))

	for i, chunk := range chunks {
		if i == 0 {
			result = append(result, chunk)
			continue
		}

		// Get overlap from previous chunk
		prevChunk := chunks[i-1]
		overlap := r.getOverlapText(prevChunk, r.ChunkOverlap)

		// Combine overlap with current chunk
		if overlap != "" {
			chunk = overlap + " " + chunk
		}

		result = append(result, chunk)
	}

	return result
}

// getOverlapText extracts the last N characters/words for overlap
func (r *RecursiveCharacterTextSplitter) getOverlapText(text string, overlapSize int) string {
	if len(text) <= overlapSize {
		return text
	}

	// Try to break at word boundaries for better overlap
	overlap := text[len(text)-overlapSize:]

	// Find the first space to avoid cutting words in half
	firstSpace := strings.IndexFunc(overlap, unicode.IsSpace)
	if firstSpace > 0 && firstSpace < len(overlap)-1 {
		overlap = overlap[firstSpace+1:]
	}

	return strings.TrimSpace(overlap)
}

// GetChunkCount returns the number of chunks that would be created
func (r *RecursiveCharacterTextSplitter) GetChunkCount(text string) int {
	return len(r.SplitText(text))
}

// GetEstimatedTokenCount estimates the number of tokens (rough approximation)
func (r *RecursiveCharacterTextSplitter) GetEstimatedTokenCount(text string) int {
	// Rough estimation: ~4 characters per token for English text
	return len(text) / 4
}

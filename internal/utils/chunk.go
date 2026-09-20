package utils

import (
	"math"
	"strings"
)

// ChunkByParagraph1
// ChunkBySentence .?!
// ChunkByFixedSize2
// ChunkByFixedSizeWithOverlap3

func ChunkByParagraph(text string) []string {
	return strings.Split(text, "\n\n")
}

func ChunkByFixedSize(text string, size int) []string {
	var chunks []string
	textLen := len(text)
	var chunk string
	var newText string = text

	loopCount := math.Ceil(float64(textLen) / float64(size))
	if loopCount < 1 {
		chunks = append(chunks, text)
		return chunks
	} else {
		for c := 0; c <= int(loopCount); c++ {
			if len(newText) < size {
				chunks = append(chunks, newText)
				return chunks
			}
			chunk = newText[0:size]
			chunks = append(chunks, chunk)
			newText = newText[size:]
			if newText == "" {
				return chunks
			}
		}
		return chunks
	}
}

func ChunkByFixedSizeWithOverlap(text string, size, overlap int) []string {
	var chunks []string
	textLen := len(text)
	var chunk string
	var newText string = text

	loopCount := math.Ceil(float64(textLen) / float64(size))
	if loopCount < 1 {
		chunks = append(chunks, text)
		return chunks
	} else {
		for len(newText) != 0 {
			if len(newText) < size {
				chunks = append(chunks, newText)
				return chunks
			}
			chunk = newText[0:size]
			chunks = append(chunks, chunk)
			newText = newText[size-overlap:]
			if len(newText) == 0 {
				return chunks
			}
		}
		return chunks
	}
}

func ChunkBySentence(text string) []string {
	var chunks []string
	var newStartIdx int = 0

	for i := 0; i < len(text); i++ {
		if text[i] == '.' || text[i] == '?' || text[i] == '!' {
			if i != len(text)-1 {
				chunk := text[newStartIdx : i+1]
				chunks = append(chunks, strings.TrimSpace(chunk))
				newStartIdx = i + 1
			} else {
				chunks = append(chunks, strings.TrimSpace(text[newStartIdx:]))
				return chunks
			}
		} else if (text[i] != '.' && text[i] != '?' && text[i] != '!') && i == len(text)-1 {
			chunks = append(chunks, strings.TrimSpace(text[newStartIdx:]))
			return chunks
		}
	}

	return chunks
}

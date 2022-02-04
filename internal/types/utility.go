package types

import (
	"math"
	"time"
)

func ChunkEvery(cards []*Card, n int) [][]*Card {
	totalChunks := int(math.Ceil(float64(len(cards)) / float64(n)))
	excess := len(cards) - ((totalChunks - 1) * n)
	chunks := make([][]*Card, totalChunks)

	for i := 0; i < totalChunks; i++ {
		if i == totalChunks-1 {
			// This is the last chunk that might have excess
			chunks[i] = make([]*Card, excess)
		} else {
			chunks[i] = make([]*Card, n)
		}
	}

	for r := 0; r < totalChunks; r++ {
		for c := 0; c < n; c++ {
			if r*n+c < len(cards) {
				chunks[r][c] = cards[r*n+c]
			}
		}
	}

	return chunks
}

func FormatTime(time time.Time) string {
	timestamp := time.Format("2006-02-01-15-04-05")
	return timestamp
}

package baleen

import (
	"math"
	"time"
)

func ChunkEvery(cards []Card, n int) [][]Card {
	totalChunks := int(math.Ceil(float64(len(cards)) / float64(n)))
	chunks := make([][]Card, totalChunks)

	for i := range chunks {
		chunks[i] = make([]Card, n)
	}

	for r := 0; r < totalChunks; r++ {
		for c := 0; c < n; c++ {
			if r*10+c < len(cards) {
				chunks[r][c] = cards[r*10+c]
			}
		}
	}

	return chunks
}

func CreateTimestamp() string {
	now := time.Now()
	timestamp := now.Format("2006-02-01-15-04-05")
	return timestamp
}

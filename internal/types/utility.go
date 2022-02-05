package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
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

// TODO: Allow users to customize the entire path
func SaveCards(cards []*Card, subfolderName string) string {
	file, _ := json.MarshalIndent(cards, "", "  ")

	folderPath := path.Join("data", subfolderName)
	exportPath := path.Join(folderPath, fmt.Sprintf("%s.json", FormatTime(time.Now())))

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Printf("Creating folder %s\n", folderPath)
		err = os.Mkdir(folderPath, 0777)
		if err != nil {

			log.Fatalf("Failed to create save folder: %v\n", err)
		}
	}

	err := ioutil.WriteFile(exportPath, file, 0644)
	if err != nil {
		log.Fatalf("Failed to save file to %s: %v\n", exportPath, err)
	}

	log.Printf("Exported to %s\n", exportPath)

	return exportPath
}

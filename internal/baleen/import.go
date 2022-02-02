package baleen

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// TODO: Grab the mapping and create database pages if not present
// TODO: Ignore any attachments that start with: https://docs.google.com/viewer?embedded=true since that's to view the page
// TODO: Add support for image attachments (for now, omit them as they are rarely useful)
func ImportToNotion(exportPath string) {
	jsonFile, err := os.Open(exportPath)
	if err != nil {
		log.Fatalf("Error occurred: %v\n", err)
	}

	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	var cards []Card
	json.Unmarshal(data, &cards)

	special, normal := 0, 0
	for _, card := range cards {
		if card.IsSpecial {
			special++
		} else {
			normal++
		}
	}

	log.Printf("special: %d, normal: %d, total: %d\n", special, normal, special+normal)
}

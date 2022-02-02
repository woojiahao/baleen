package baleen

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/woojiahao/baleen/internal/config"
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

	config := config.New("configs/conf.json")
	prepareDatabases(config)
}

func prepareDatabases(config *config.Config) {
	databaseNames := config.DatabaseNames()

	for _, name := range databaseNames {

	}
}

function getDatabaseIdsFromNames(names []string) []string {
	var ids []string

	for _, name := range names {

	}
}
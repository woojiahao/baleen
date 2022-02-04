package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	TrelloKey   string
	TrelloToken string
	NotionKey   string
}

func New(envPath string) *Env {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Failed to load environment from file %s: %v\n", envPath, err)
	}

	return &Env{
		os.Getenv("TRELLO_API_KEY"),
		os.Getenv("TRELLO_TOKEN"),
		os.Getenv("NOTION_INTEGRATION_KEY"),
	}
}

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/woojiahao/baleen/internal/baleen"
)

// TODO: Support general migrations from Trello to Notion
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	baleen.ExportTrelloBoard("Programming Bucket")
}

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

	// cards := baleen.ExtractTrelloBoard("Programming Bucket")
	// baleen.SaveCards(cards)
	cards := baleen.LoadCardsFromExport("data/2022-02-02-20-27-35-trello.json")
	baleen.ImportToNotion(cards)
}

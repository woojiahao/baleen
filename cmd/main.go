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

	cards := baleen.ExtractTrelloBoard("Programming Bucket")
	baleen.SaveCards(cards)

	// notion.UpdateDatabaseProperties(
	// 	"d583efbe-a96d-49ca-afc5-9d7566c051da",
	// 	notion.NotionProperties{
	// 		"diff": notion.NotionProperty{
	// 			string(notion.RichText): notion.NotionPropertyBody{
	// 				MultiSelectOptions: nil,
	// 			},
	// 		},
	// 		"multi": notion.NotionProperty{
	// 			string(notion.MultiSelect): notion.NotionPropertyBody{
	// 				MultiSelectOptions: []notion.NotionMultiSelect{
	// 					{Name: "Something", Color: "blue"},
	// 					{Name: "Something Else", Color: "green"},
	// 				},
	// 			},
	// 		},
	// 	},
	// )
}

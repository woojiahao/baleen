package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/woojiahao/baleen/internal/api/trello"
)

// TODO: Allow users to get the board ID of any board by searching by name
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	programmingBucketId := trello.FindBoardIdByName("Programming Bucket")
	log.Printf("Id: %s\n", programmingBucketId)
	// trello.GetLists()

	fmt.Println(os.Getenv("TRELLO_API_KEY"))
}

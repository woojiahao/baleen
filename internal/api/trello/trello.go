package trello

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/woojiahao/baleen/internal/api"
)

const (
	trelloApi = "https://api.trello.com/1"
)

type SearchBoardsResult struct {
	Boards []struct {
		Id   string
		Name string
	}
}

func baseQueryParams() map[string]string {
	_ = godotenv.Load(".env")
	return map[string]string{
		"key":   os.Getenv("TRELLO_API_KEY"),
		"token": os.Getenv("TRELLO_TOKEN"),
	}
}

func GetBoards() {
	api.Get(trelloApi, "/members/me/boards", baseQueryParams())
}

func FindBoardIdByName(name string) string {
	queryParams := baseQueryParams()
	queryParams["modelTypes"] = "boards"
	queryParams["query"] = name
	resp := api.Get(trelloApi, "/search", queryParams)

	var results SearchBoardsResult
	json.Unmarshal(resp, &results)

	return results.Boards[0].Id
}

func GetLists() {
	queryParams := baseQueryParams()
	boardId := os.Getenv("TRELLO_BOARD_ID")
	queryParams["filter"] = "open"
	api.Get(trelloApi, fmt.Sprintf("/boards/%s/lists", boardId), queryParams)
}

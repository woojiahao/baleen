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

	var results TrelloSearchBoards
	json.Unmarshal(resp, &results)

	return results.Boards[0].Id
}

func GetListsInBoard(boardId string) []TrelloBoardLists {
	queryParams := baseQueryParams()
	queryParams["filter"] = "open"
	resp := api.Get(trelloApi, fmt.Sprintf("/boards/%s/lists", boardId), queryParams)

	var results []TrelloBoardLists
	json.Unmarshal(resp, &results)

	return results
}

func GetCardsInList(listId string) []TrelloBasicCard {
	resp := api.Get(trelloApi, fmt.Sprintf("/lists/%s/cards", listId), baseQueryParams())

	var cards []TrelloBasicCard
	json.Unmarshal(resp, &cards)

	return cards
}

func GetFullCard(cardId string) TrelloFullCard {
	params := baseQueryParams()
	params["actions"] = "commentCard"
	params["attachments"] = "true"
	params["fields"] = "name"
	params["attachment_fields"] = "all"

	resp := api.Get(trelloApi, fmt.Sprintf("/cards/%s", cardId), params)

	var card TrelloFullCard
	json.Unmarshal(resp, &card)

	return card
}

package baleen

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"

	"github.com/woojiahao/baleen/internal/api/trello"
)

// For every card that is found, we will want to extract any comments and attachments found
func ExportTrelloBoard(boardName string) {
	programmingBucketId := trello.FindBoardIdByName(boardName)
	lists := trello.GetListsInBoard(programmingBucketId)

	var cards []Card

	// Process every single card and extract the information we want to export
	for _, list := range lists {
		for _, card := range trello.GetCardsInList(list.Id) {
			if card.Badges.Attachments > 0 || card.Badges.Comments > 0 {
				var labels []Label

				for _, label := range card.Labels {
					labels = append(labels, Label{label.Name, label.Color})
				}

				baleenCard := Card{
					Id:             card.Id,
					Name:           card.Name,
					ParentListName: list.Name,
					Labels:         labels,
					LastUpdate:     card.LastUpdate,
					IsSpecial:      (card.Badges.Attachments > 0 || card.Badges.Comments > 0),
					Comments:       make([]string, 0),
					Attachments:    make([]Attachment, 0),
				}

				cards = append(cards, baleenCard)
			}
		}
	}

	specialCards := processSpecialCards(cards)
	log.Printf("special cards: %v\n", specialCards)

	var normalCards []Card
	for _, card := range cards {
		if !card.IsSpecial {
			normalCards = append(normalCards, card)
		}
	}

	all_cards := append(normalCards, specialCards...)
	file, _ := json.MarshalIndent(all_cards, "", " ")
	_ = ioutil.WriteFile("data/trello.json", file, 0644)
}

// To speed up the processing of special cards, we will attempt to parallelize the API calls
// Given the rate limit of 100 requests/10s or 10 requests/s, we will dispatch 10 goroutines at a time to process 10
// cards simultaneously
func processSpecialCards(cards []Card) []Card {
	var specialCards []Card
	for _, card := range cards {
		if card.IsSpecial {
			specialCards = append(specialCards, card)
		}
	}

	// TODO: Unit test this to ensure that it's actually running in parallel
	chunks := chunkEvery(specialCards, 10)
	var updatedSpecials []Card
	for _, chunk := range chunks {
		c := make(chan []Card)
		go retrieveSpecial(chunk, c)
		fullSpecials := <-c
		updatedSpecials = append(updatedSpecials, fullSpecials...)
	}

	return updatedSpecials
}

func retrieveSpecial(cards []Card, c chan []Card) {
	newCards := make([]Card, len(cards))
	for _, card := range cards {
		full_card := trello.GetFullCard(card.Id)

		// Load comments
		for _, comment := range full_card.Actions {
			card.Comments = append(card.Comments, comment.Data.Text)
		}

		// Load attachments
		for _, attachment := range full_card.Attachments {
			card.Attachments = append(card.Attachments, Attachment{
				IsUpload: attachment.IsUpload,
				Name:     attachment.Name,
				Url:      attachment.Url,
				Filename: attachment.Filename,
			})
		}

		newCards = append(newCards, card)
	}

	c <- newCards
}

func chunkEvery(cards []Card, n int) [][]Card {
	totalChunks := int(math.Ceil(float64(len(cards)) / float64(n)))
	chunks := make([][]Card, totalChunks)

	for i := range chunks {
		chunks[i] = make([]Card, n)
	}

	for r := 0; r < totalChunks; r++ {
		for c := 0; c < n; c++ {
			if r*10+c < len(cards) {
				chunks[r][c] = cards[r*10+c]
			}
		}
	}

	return chunks
}

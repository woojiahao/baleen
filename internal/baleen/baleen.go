package baleen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

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
			var labels []Label

			for _, label := range card.Labels {
				labels = append(labels, Label{label.Name, label.Color})
			}

			cards = append(cards, Card{
				Id:             card.Id,
				Name:           card.Name,
				ParentListName: list.Name,
				Labels:         labels,
				LastUpdate:     card.LastUpdate,
				IsSpecial:      (card.Badges.Attachments > 0 || card.Badges.Comments > 0),
				Comments:       make([]string, 0),
				Attachments:    make([]Attachment, 0),
			})
		}
	}

	specialCards := processSpecialCards(cards)

	var normalCards []Card
	for _, card := range cards {
		if !card.IsSpecial {
			normalCards = append(normalCards, card)
		}
	}

	var allCards []Card
	allCards = append(allCards, normalCards...)
	allCards = append(allCards, specialCards...)

	file, _ := json.MarshalIndent(allCards, "", " ")
	exportPath := path.Join("data", fmt.Sprintf("%s-trello.json", CreateTimestamp()))
	_ = ioutil.WriteFile(exportPath, file, 0644)
	log.Printf("Export to %s\n", exportPath)
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
	chunks := ChunkEvery(specialCards, 10)
	var updatedSpecials []Card
	for _, chunk := range chunks {
		c := make(chan []Card)
		go retrieveSpecial(chunk[:len(chunk)/2], c)
		go retrieveSpecial(chunk[len(chunk)/2:], c)

		fullSpecialsLeft, fullSpecialsRight := <-c, <-c
		updatedSpecials = append(updatedSpecials, fullSpecialsLeft...)
		updatedSpecials = append(updatedSpecials, fullSpecialsRight...)
	}

	return updatedSpecials
}

func retrieveSpecial(cards []Card, c chan []Card) {
	var newCards []Card
	for _, card := range cards {
		fullCard := trello.GetFullCard(card.Id)

		// Load comments
		for _, comment := range fullCard.Actions {
			card.Comments = append(card.Comments, comment.Data.Text)
		}

		// Load attachments
		for _, attachment := range fullCard.Attachments {
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

package trello

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"

	t "github.com/adlio/trello"
	"github.com/woojiahao/baleen/internal/env"
	"github.com/woojiahao/baleen/internal/types"
)

func ExtractTrelloBoard(boardName string) []*types.Card {
	log.Printf("Extracting Trello board %s\n", boardName)

	env := env.New(".env")
	client := t.NewClient(env.TrelloKey, env.TrelloToken)

	boards, _ := client.SearchBoards(boardName)
	lists, _ := boards[0].GetLists()

	var normalCards, specialCards []*types.Card

	for _, list := range lists {
		log.Printf("Extracting %s\n", list.Name)

		cards, _ := list.GetCards()

		for _, card := range cards {
			var labels []*types.Label
			for _, label := range card.Labels {
				labels = append(labels, &types.Label{
					Name:  label.Name,
					Color: label.Color,
				})
			}

			isSpecial := card.Badges.Attachments > 0 || card.Badges.Comments > 0

			typesCard := &types.Card{
				Id:             card.ID,
				Name:           card.Name,
				Description:    card.Desc,
				ParentListName: list.Name,
				Labels:         labels,
				LastUpdate:     card.DateLastActivity,
				IsSpecial:      isSpecial,
				Comments:       []string{},
				Attachments:    []*types.Attachment{},
			}

			if isSpecial {
				specialCards = append(specialCards, typesCard)
			} else {
				normalCards = append(normalCards, typesCard)
			}
		}
	}

	specialCards = processSpecialCards(client, specialCards)

	var typesCards []*types.Card
	typesCards = append(typesCards, normalCards...)
	typesCards = append(typesCards, specialCards...)

	return typesCards
}

func SaveCards(cards []*types.Card) string {
	file, _ := json.MarshalIndent(cards, "", "  ")
	exportPath := path.Join("data", fmt.Sprintf("%s-trello.json", types.FormatTime(time.Now())))
	_ = ioutil.WriteFile(exportPath, file, 0644)
	log.Printf("Export to %s\n", exportPath)

	return exportPath
}

func processSpecialCards(client *t.Client, specialCards []*types.Card) []*types.Card {
	chunks := types.ChunkEvery(specialCards, 10)

	var updatedSpecials []*types.Card
	for i, chunk := range chunks {
		c := make(chan []*types.Card)
		go parallelProcessSpecial(client, chunk[:len(chunk)/2], c)
		go parallelProcessSpecial(client, chunk[len(chunk)/2:], c)

		specialLeft, specialRight := <-c, <-c
		updatedSpecials = append(updatedSpecials, specialLeft...)
		updatedSpecials = append(updatedSpecials, specialRight...)

		log.Printf("Processes %d/%d\n", i+1, len(chunks))
	}

	return updatedSpecials
}

func parallelProcessSpecial(client *t.Client, specialCards []*types.Card, c chan []*types.Card) {
	var cards []*types.Card
	for _, card := range specialCards {
		comments, attachments := getSpecial(client, card.Id)
		card.Comments = comments
		card.Attachments = attachments
		cards = append(cards, card)
	}

	c <- cards
}

func getSpecial(client *t.Client, cardId string) ([]string, []*types.Attachment) {
	comments := []string{}
	attachments := []*types.Attachment{}

	var specialCard *t.Card
	client.Get(
		fmt.Sprintf("cards/%s", cardId),
		map[string]string{
			"actions":           "commentCard",
			"attachments":       "true",
			"fields":            "name",
			"attachment_fields": "all",
		},
		&specialCard,
	)

	for _, action := range specialCard.Actions {
		if action.Type == "commentCard" {
			comments = append(comments, action.Data.Text)
		}
	}

	for _, attachment := range specialCard.Attachments {
		attachments = append(attachments, &types.Attachment{
			IsUpload: attachment.IsUpload,
			Name:     attachment.Name,
			Url:      attachment.URL,
		})
	}

	return comments, attachments
}

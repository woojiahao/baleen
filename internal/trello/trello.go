package trello

import (
	"fmt"
	"log"
	"path"

	t "github.com/adlio/trello"
	"github.com/woojiahao/baleen/internal/env"
	"github.com/woojiahao/baleen/internal/types"
)

func ArchiveAll(boardName, envPath string) {
	log.Printf("Archiving all lists in %s\n", boardName)

	env := env.New(envPath)
	client := t.NewClient(env.TrelloKey, env.TrelloToken)
	lists := getLists(client, boardName)

	for _, list := range lists {
		log.Printf("Archiving %s\n", list.Name)

		var burner interface{}

		client.Post(
			path.Join("lists", list.ID, "archiveAllCards"),
			t.Arguments{},
			burner,
		)

		log.Printf("Archived %s\n", list.Name)
	}
}

func ExportTrelloBoard(boardName, envPath string) []*types.Card {
	log.Printf("Extracting Trello board %s\n", boardName)

	env := env.New(envPath)
	client := t.NewClient(env.TrelloKey, env.TrelloToken)
	lists := getLists(client, boardName)

	var normalCards, specialCards []*types.Card

	for _, list := range lists {
		log.Printf("Extracting %s\n", list.Name)

		cards, err := list.GetCards()
		if err != nil {
			log.Fatalf("Failed to get cards from %s: %v\n", boardName, err)
		}

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
	typesCards = append(typesCards, specialCards...)
	typesCards = append(typesCards, normalCards...)

	return typesCards
}

func getLists(client *t.Client, boardName string) []*t.List {
	boards, err := client.SearchBoards(boardName)
	if err != nil {
		log.Fatalf("Failed to find %s: %v\n", boardName, err)
	}

	lists, err := boards[0].GetLists()
	if err != nil {
		log.Fatalf("Failed to get lists of board %s: %v\n", boardName, err)
	}

	return lists
}

func processSpecialCards(client *t.Client, specialCards []*types.Card) []*types.Card {
	log.Println("Processing special cards...")

	chunks := types.ChunkEvery(specialCards, 10)

	var updatedSpecials []*types.Card
	c := make(chan []*types.Card)

	for i, chunk := range chunks {
		go parallelProcessSpecial(client, chunk[:len(chunk)/2], c)
		go parallelProcessSpecial(client, chunk[len(chunk)/2:], c)

		specialLeft, specialRight := <-c, <-c
		updatedSpecials = append(updatedSpecials, specialLeft...)
		updatedSpecials = append(updatedSpecials, specialRight...)

		if (i+1)%10 == 0 {
			log.Printf("Processed %d/%d\n", i+1, len(chunks))
		}
	}

	log.Printf("Processed all special cards!")

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

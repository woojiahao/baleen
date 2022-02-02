package baleen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/adlio/trello"
	"github.com/joho/godotenv"
)

// TODO: Parallelize extraction of special data
func ExtractTrelloBoard(boardName string) []Card {
	log.Printf("Extracting Trello board %s\n", boardName)

	_ = godotenv.Load(".env")
	client := trello.NewClient(os.Getenv("TRELLO_API_KEY"), os.Getenv("TRELLO_TOKEN"))
	boards, _ := client.SearchBoards(boardName)
	lists, _ := boards[0].GetLists()

	var normalCards, specialCards []Card

	for _, list := range lists {
		log.Printf("Extracting %s\n", list.Name)

		cards, _ := list.GetCards()

		for _, card := range cards {
			var labels []Label
			for _, label := range card.Labels {
				labels = append(labels, Label{label.Name, label.Color})
			}

			isSpecial := card.Badges.Attachments > 0 || card.Badges.Comments > 0

			baleenCard := Card{
				Id:             card.ID,
				Name:           card.Name,
				Description:    card.Desc,
				ParentListName: list.Name,
				Labels:         labels,
				LastUpdate:     FormatTime(*card.DateLastActivity),
				IsSpecial:      isSpecial,
				Comments:       []string{},
				Attachments:    []Attachment{},
			}

			if isSpecial {
				specialCards = append(specialCards, baleenCard)
			} else {
				normalCards = append(normalCards, baleenCard)
			}
		}
	}

	specialCards = processSpecialCards(client, specialCards)

	var baleenCards []Card
	baleenCards = append(baleenCards, normalCards...)
	baleenCards = append(baleenCards, specialCards...)

	return baleenCards
}

func SaveCards(cards []Card) string {
	file, _ := json.MarshalIndent(cards, "", "  ")
	exportPath := path.Join("data", fmt.Sprintf("%s-trello.json", FormatTime(time.Now())))
	_ = ioutil.WriteFile(exportPath, file, 0644)
	log.Printf("Export to %s\n", exportPath)

	return exportPath
}

func processSpecialCards(client *trello.Client, specialCards []Card) []Card {
	chunks := ChunkEvery(specialCards, 10)

	var updatedSpecials []Card
	for _, chunk := range chunks {
		c := make(chan []Card)
		go parallelProcessSpecial(client, chunk[:len(chunk)/2], c)
		go parallelProcessSpecial(client, chunk[len(chunk)/2:], c)

		specialLeft, specialRight := <-c, <-c
		updatedSpecials = append(updatedSpecials, specialLeft...)
		updatedSpecials = append(updatedSpecials, specialRight...)
	}

	return updatedSpecials
}

func parallelProcessSpecial(client *trello.Client, specialCards []Card, c chan []Card) {
	var cards []Card
	var ids []string
	for _, card := range specialCards {
		comments, attachments := getSpecial(client, card.Id)
		card.Comments = comments
		card.Attachments = attachments
		cards = append(cards, card)
		ids = append(ids, card.Id)
	}

	log.Printf("Processed %v\n", ids)

	c <- cards
}

func getSpecial(client *trello.Client, cardId string) ([]string, []Attachment) {
	comments := []string{}
	attachments := []Attachment{}

	var specialCard *trello.Card
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
		attachments = append(attachments, Attachment{
			IsUpload: false,
			Name:     attachment.Name,
			Url:      attachment.URL,
		})
	}

	return comments, attachments
}

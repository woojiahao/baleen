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

	var baleenCards []Card

	for _, list := range lists {
		log.Printf("Extracting %s\n", list.Name)

		cards, _ := list.GetCards()

		for _, card := range cards {
			var labels []Label
			for _, label := range card.Labels {
				labels = append(labels, Label{label.Name, label.Color})
			}

			comments, attachments := getSpecial(client, card)

			baleenCard := Card{
				Id:             card.ID,
				Name:           card.Name,
				Description:    card.Desc,
				ParentListName: list.Name,
				Labels:         labels,
				LastUpdate:     FormatTime(*card.DateLastActivity),
				IsSpecial:      card.Badges.Attachments > 0 || card.Badges.Comments > 0,
				Comments:       comments,
				Attachments:    attachments,
			}

			baleenCards = append(baleenCards, baleenCard)
		}
	}

	return baleenCards
}

func SaveCards(cards []Card) string {
	file, _ := json.MarshalIndent(cards, "", "  ")
	exportPath := path.Join("data", fmt.Sprintf("%s-trello.json", FormatTime(time.Now())))
	_ = ioutil.WriteFile(exportPath, file, 0644)
	log.Printf("Export to %s\n", exportPath)

	return exportPath
}

func getSpecial(client *trello.Client, card *trello.Card) (comments []string, attachments []Attachment) {
	comments, attachments = []string{}, []Attachment{}

	if !(card.Badges.Attachments > 0 || card.Badges.Comments > 0) {
		return
	}

	var specialCard *trello.Card
	client.Get(
		fmt.Sprintf("cards/%s", card.ID),
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

	return
}

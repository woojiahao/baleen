package baleen

import (
	"github.com/woojiahao/baleen/internal/notion"
	"github.com/woojiahao/baleen/internal/trello"
)

func Migrate(trelloBoardName, configPath, envPath string, toSave bool) {
	cards := trello.ExtractTrelloBoard(trelloBoardName)
	if toSave {
		trello.SaveCards(cards)
	}

	notion.ImportToNotion(cards)
}

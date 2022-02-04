package baleen

import (
	"github.com/woojiahao/baleen/internal/notion"
	"github.com/woojiahao/baleen/internal/trello"
)

// Performs full migration from Trello board to Notion
func Migrate(trelloBoardName, configPath, envPath string, toSave bool) {
	cards := trello.ExportTrelloBoard(trelloBoardName, envPath)
	if toSave {
		trello.SaveCards(cards)
	}

	notion.ImportToNotion(cards, envPath, configPath)
}

// Imports into Notion from existing save file
func Import(savePath, configPath, envPath string) {
	cards := notion.LoadSave(savePath)
	notion.ImportToNotion(cards, envPath, configPath)
}

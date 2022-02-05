package baleen

import (
	"github.com/woojiahao/baleen/internal/notion"
	"github.com/woojiahao/baleen/internal/trello"
	"github.com/woojiahao/baleen/internal/types"
)

const (
	savePath = "saves"
)

// Performs full migration from Trello board to Notion
func Migrate(trelloBoardName, configPath, envPath string, toSave bool) {
	cards := trello.ExportTrelloBoard(trelloBoardName, envPath)
	if toSave {
		types.SaveCards(cards, savePath)
	}

	notion.ImportToNotion(cards, envPath, configPath)
}

// Imports into Notion from existing save file
func Import(savePath, configPath, envPath string) {
	cards := notion.LoadSave(savePath)
	notion.ImportToNotion(cards, envPath, configPath)
}

func ExportAndSave(trelloBoardName, envPath string) {
	cards := trello.ExportTrelloBoard(trelloBoardName, envPath)
	types.SaveCards(cards, savePath)
}

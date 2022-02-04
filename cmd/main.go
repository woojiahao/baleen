package main

import "github.com/woojiahao/baleen/internal/baleen"

// TODO: Support general migrations from Trello to Notion
func main() {
	// baleen.Migrate("Programming Bucket", "configs/conf.json", ".env", true)
	baleen.Import("data/saves/2022-04-02-22-18-01.json", "configs/conf.json", ".env")
}

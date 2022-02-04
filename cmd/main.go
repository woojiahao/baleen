package main

import "github.com/woojiahao/baleen/internal/baleen"

// TODO: Support general migrations from Trello to Notion
func main() {
	baleen.Migrate("Programming Bucket", "configs/conf.json", ".env", true)
	// baleen.Import("")
}

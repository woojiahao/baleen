package notion

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/woojiahao/baleen/internal/api"
)

const (
	notionApi     = "https://api.notion.com/v1"
	notionVersion = "2021-08-16"
)

func headers() map[string]string {
	_ = godotenv.Load(".env")
	return map[string]string{
		"Authorization":  fmt.Sprintf("Bearer %s", os.Getenv("NOTION_INTEGRATION_KEY")),
		"Content-Type":   "application/json",
		"Notion-Version": notionVersion,
	}
}

func GetAllPages() {
	resp := api.Post(notionApi, "/search", headers(), map[string]string{
		"query": "",
	})
	log.Printf("all pages: %v\n", string(resp))
}

// TODO: Ensure that the page is a database for import
func GetPageByName(pageName string) {
	resp := api.Post(notionApi, "/search", headers(), map[string]string{
		"query": pageName,
	})
	log.Printf("page: %v\n", string(resp))
}

func UpdateDatabaseProperties(databaseId string, properties NotionProperties) {
	log.Printf("properties: %v\n", properties)
	propertiesJson, _ := json.Marshal(&properties)
	log.Printf("propertiesJson: %v\n", string(propertiesJson))
	body := map[string]string{
		"properties": string(propertiesJson),
	}
	resp := api.Patch(notionApi, fmt.Sprintf("/databases/%s", databaseId), headers(), body)
	log.Printf("updated: %v\n", string(resp))
}

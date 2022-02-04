package notion

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/woojiahao/baleen/internal/baleen"
	"github.com/woojiahao/baleen/internal/config"
	"github.com/woojiahao/baleen/internal/types"
	"golang.org/x/net/context"
)

// Types for database metadata
type (
	databaseName    string
	databaseId      string
	databaseNameIds map[databaseName]databaseId
)

// Types for Notion API data
type (
	primaryLink string
)

type (
	attachmentName string
	attachmentLink string
	attachments    map[attachmentName]attachmentLink
)

func (a *attachments) first() (attachmentName, attachmentLink) {
	for name, link := range *a {
		return name, link
	}

	return "null", "null"
}

func (a *attachments) toMap() map[string]string {
	m := make(map[string]string)
	for name, link := range *a {
		m[string(name)] = string(link)
	}

	return m
}

// Import a set of cards into Notion. The cards should either be loaded from a file with LoadCardsFromExport or directly
// from trello.ExtractTrelloBoard.
func ImportToNotion(cards []*types.Card) {
	log.Printf("Importing cards into Notion\n")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Failed to load .env: %v\n", err)
	}

	notionIntegrationKey := os.Getenv("NOTION_INTEGRATION_KEY")
	notion := notionapi.NewClient(notionapi.Token(notionIntegrationKey))

	config := config.New("configs/conf.json")
	nameIds := getDatabaseNameIds(notion, config.DatabaseNames())
	labels := extractLabels(config, cards)

	addDatabaseProperties(notion, nameIds, labels)
	addCardToRespectiveDatabase(config, notion, nameIds, cards)
}

func LoadCardsFromExport(exportPath string) []*types.Card {
	log.Printf("Loading cards from export path %s\n", exportPath)

	jsonFile, err := os.Open(exportPath)
	if err != nil {
		log.Fatalf("Error occurred: %v\n", err)
	}

	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	var cards []*types.Card
	json.Unmarshal(data, &cards)

	return cards
}

// Add a card to its respective database
func addCardToRespectiveDatabase(config *config.Config, notion *notionapi.Client, nameIds *databaseNameIds, cards []*Card) {
	log.Printf("Adding cards to database")

	for _, card := range cards {
		log.Printf("Adding card %s\n", card.Name)

		fileAttachments, urlAttachments := organizeAttachments(card)

		_, pl := urlAttachments.first()

		properties := createProperties(card, primaryLink(pl))
		children := createChildren(fileAttachments, urlAttachments, card.Comments)

		_, err := notion.Page.Create(context.Background(), &notionapi.PageCreateRequest{
			Parent: notionapi.Parent{
				DatabaseID: notionapi.DatabaseID((*nameIds)[databaseName(config.Database[card.ParentListName])]),
			},
			Properties: properties,
			Children:   children,
		})

		if err != nil {
			log.Fatalf("Error occurred when adding cards to database: %v\n", err)
		}
	}
}

// Add the necessary properties for importing Trello information into a Notion card
// Properties include a Description, Primary Link (first link in the attachments), Labels, and Last Updated
func addDatabaseProperties(notion *notionapi.Client, nameIds *databaseNameIds, labels []*types.Label) {
	log.Printf("Adding properties to database")

	var labelOptions []notionapi.Option
	for _, label := range labels {
		labelOptions = append(labelOptions, notionapi.Option{
			Name:  label.Name,
			Color: notionapi.Color(label.Color),
		})
	}

	for name, id := range *nameIds {
		request := &notionapi.DatabaseUpdateRequest{
			Properties: notionapi.PropertyConfigs{
				"Description":  richTextConfig(),
				"Primary Link": urlConfig(),
				"Labels":       multiSelectConfig(labelOptions),
				"Last Updated": dateConfig(),
			},
		}
		_, err := notion.Database.Update(context.Background(), notionapi.DatabaseID(id), request)

		if err != nil {
			log.Fatalf("Failed to add properties to %s: %v\n", name, err)
		}
	}
}

func createProperties(card *types.Card, pl primaryLink) notionapi.Properties {
	labelOptions := organizeLabels(card.Labels)

	properties := notionapi.Properties{}

	properties["Name"] = titleProperty(card.Name)
	properties["Description"] = richTextProperty(card.Description, noLink)

	if pl != "null" {
		properties["Primary Link"] = urlProperty(string(pl))
	}

	if card.LastUpdate != nil {
		properties["Last Updated"] = dateProperty(card.LastUpdate)
	}

	if len(labelOptions) > 0 {
		properties["Labels"] = multiSelectProperty(labelOptions)
	}

	return properties
}

func createChildren(fileAttachments, urlAttachments *attachments, comments []string) []notionapi.Block {
	var children []notionapi.Block

	children = append(children, heading1("File Attachments"))
	for _, file := range linkBlocks(fileAttachments.toMap()) {
		children = append(children, file)
	}

	children = append(children, heading1("URL Attachments"))
	for _, url := range linkBlocks(urlAttachments.toMap()) {
		children = append(children, url)
	}

	children = append(children, heading1("Comments"))
	for _, comment := range comments {
		children = append(children, paragraph(comment))
	}

	return children
}

func organizeLabels(labels []Label) []notionapi.Option {
	var options []notionapi.Option
	for _, label := range labels {
		options = append(options, notionapi.Option{
			Name: label.Name,
		})
	}
	return options
}

func organizeAttachments(card *Card) (fileAttachments, urlAttachments *attachments) {
	files, urls := make(attachments), make(attachments)

	for _, attachment := range card.Attachments {
		if attachment.IsUpload {
			files[attachmentName(attachment.Name)] = attachmentLink(attachment.Url)
		} else {
			urls[attachmentName(attachment.Name)] = attachmentLink(attachment.Url)
		}
	}

	fileAttachments, urlAttachments = &files, &urls

	return
}

func extractLabels(config *config.Config, cards []*Card) []*Label {
	labelsMap := make(map[string]string)

	for _, card := range cards {
		for _, label := range card.Labels {
			if alt, ok := config.Color[label.Color]; ok {
				labelsMap[label.Name] = alt
			} else {
				labelsMap[label.Name] = label.Color
			}
		}
	}

	var labels []*Label
	for name, color := range labelsMap {
		labels = append(labels, &Label{name, color})
	}

	return labels
}

func getDatabaseNameIds(notion *notionapi.Client, names []string) *databaseNameIds {
	nameIds := make(databaseNameIds)

	searchResp, _ := notion.Search.Do(context.Background(), &notionapi.SearchRequest{
		Filter: map[string]string{
			"value":    "database",
			"property": "object",
		},
	})

	for _, result := range searchResp.Results {
		r := result.(*notionapi.Database)
		title := r.Title[0].Text.Content
		if baleen.Contains(title, names) {
			nameIds[databaseName(title)] = databaseId(r.ID.String())
		}
	}

	for _, name := range names {
		if _, ok := nameIds[databaseName(name)]; !ok {
			log.Fatalf("Unable to find database titled %s to import to\n", name)
		}
	}

	return &nameIds
}

package notion

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/jomei/notionapi"
	"github.com/woojiahao/baleen/internal/config"
	"github.com/woojiahao/baleen/internal/env"
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
func ImportToNotion(cards []*types.Card, envPath, configPath string) {
	log.Printf("Importing cards into Notion\n")

	env := env.New(envPath)
	notion := notionapi.NewClient(notionapi.Token(env.NotionKey))

	config := config.New(configPath)
	nameIds := getDatabaseNameIds(notion, config.DatabaseNames())
	labels := extractLabels(config, cards)

	addDatabaseProperties(config, notion, nameIds, labels)
	importCards(config, notion, nameIds, cards)
}

func LoadSave(exportPath string) []*types.Card {
	log.Printf("Loading cards from save %s\n", exportPath)

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
func importCards(config *config.Config, notion *notionapi.Client, nameIds *databaseNameIds, cards []*types.Card) {
	log.Printf("Adding cards to database\n")

	chunks := types.ChunkEvery(cards, 3)
	c := make(chan *types.Card, 3)
	var errCards []*types.Card

	for i, chunk := range chunks {
		for _, card := range chunk {
			go importCard(config, notion, nameIds, card, c)
		}

		f, s, t := <-c, <-c, <-c
		if f != nil {
			errCards = append(errCards, f)
		}

		if s != nil {
			errCards = append(errCards, s)
		}

		if t != nil {
			errCards = append(errCards, t)
		}

		if (i+1)%15 == 0 {
			log.Printf("Imported %d/%d\n", i+1, len(chunks))
		}
	}

	log.Printf("Imported all cards!\n")

	if errCards != nil {
		errPath := types.SaveCards(errCards, "errors")
		log.Printf("Saved error cards to %s\n", errPath)
	}
}

// Emits nil to the channel if the card didn't occur any errors. A returned card will indicate that the card has failed
// to be added.
func importCard(
	config *config.Config,
	notion *notionapi.Client,
	nameIds *databaseNameIds,
	card *types.Card,
	c chan *types.Card,
) {
	fileAttachments, urlAttachments := organizeAttachments(card)

	_, pl := urlAttachments.first()

	properties := createProperties(config, card, primaryLink(pl))
	children := createChildren(fileAttachments, urlAttachments, card.Description, card.Comments)

	tries := 1

	request := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID((*nameIds)[databaseName(config.Database[card.ParentListName])]),
		},
		Properties: properties,
		Children:   children,
	}

	_, err := notion.Page.Create(context.Background(), request)

	// Try three more times in case it's just a timeout issue
	for err != nil && tries < 4 {
		tries++
		log.Printf("Error occurred when adding card (%s) to database: %v\n", card.Name, err)
		j, _ := json.MarshalIndent(request, "", "  ")
		log.Printf("Request was: %v\n", string(j))
		log.Printf("Trying again after 2s wait... (Try %d)\n", tries)
		time.Sleep(2 * time.Second)
		_, err = notion.Page.Create(context.Background(), request)
	}

	if tries >= 4 {
		log.Printf("Unable to add card %s. Saving for inspection later.\n", card.Name)
		c <- card
	} else {
		log.Printf("Added card %s\n", card.Name)
		c <- nil
	}
}

// Add the necessary properties for importing Trello information into a Notion card
// Properties include a Description, Primary Link (first link in the attachments), Labels, and Last Updated
func addDatabaseProperties(
	config *config.Config,
	notion *notionapi.Client,
	nameIds *databaseNameIds,
	labels []*types.Label,
) {
	log.Printf("Adding properties to database")

	labelOptions := organizeLabels(config, labels)

	for name, id := range *nameIds {
		properties := make(notionapi.PropertyConfigs)
		properties["Description"] = richTextConfig()
		properties["Primary Link"] = urlConfig()
		if len(labelOptions) > 0 {
			properties["Labels"] = multiSelectConfig(labelOptions)
		}
		properties["Last Updated"] = dateConfig()

		request := &notionapi.DatabaseUpdateRequest{Properties: properties}
		_, err := notion.Database.Update(context.Background(), notionapi.DatabaseID(id), request)

		if err != nil {
			log.Fatalf("Failed to add properties to %s: %v\n", name, err)
		}
	}
}

func createProperties(config *config.Config, card *types.Card, pl primaryLink) notionapi.Properties {
	labelOptions := organizeLabels(config, card.Labels)

	properties := notionapi.Properties{}

	properties["Name"] = titleProperty(card.Name)

	description := card.Description
	descriptionLimit := 100
	if len(description) > descriptionLimit {
		description = description[:descriptionLimit]
	}
	properties["Description"] = richTextProperty(description, noLink)

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

func createChildren(fileAttachments, urlAttachments *attachments, description string, comments []string) []notionapi.Block {
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

	children = append(children, heading1("Description"))
	children = append(children, paragraph(description))

	return children
}

// Organize a list of labels into a list of options for Notion
func organizeLabels(config *config.Config, labels []*types.Label) []notionapi.Option {
	if labels == nil {
		return []notionapi.Option{}
	}

	var options []notionapi.Option
	for _, label := range labels {
		option := notionapi.Option{
			Name:  label.Name,
			Color: notionapi.Color(label.Color),
		}
		if alt, ok := config.Color[label.Color]; ok {
			option.Color = notionapi.Color(alt)
		}

		options = append(options, option)
	}
	return options
}

func organizeAttachments(card *types.Card) (fileAttachments, urlAttachments *attachments) {
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

// Extracts the unique labels from a list of cards
func extractLabels(config *config.Config, cards []*types.Card) []*types.Label {
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

	var labels []*types.Label
	for name, color := range labelsMap {
		labels = append(labels, &types.Label{name, color})
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
		if contains(title, names) {
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

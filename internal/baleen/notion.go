package baleen

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/woojiahao/baleen/internal/config"
	"golang.org/x/net/context"
)

// TODO: Ignore any attachments that start with: https://docs.google.com/viewer?embedded=true since that's to view the page
// TODO: Add support for image attachments (for now, omit them as they are rarely useful)
// TODO: Clean up the high nested structure
func ImportToNotion(cards []Card) {
	log.Printf("Importing cards into Notion\n")

	_ = godotenv.Load(".env")
	notionIntegrationKey := os.Getenv("NOTION_INTEGRATION_KEY")
	notion := notionapi.NewClient(notionapi.Token(notionIntegrationKey))

	config := config.New("configs/conf.json")
	ids := getDatabaseIdsFromNames(notion, config.DatabaseNames())
	labels := extractLabels(config, cards)
	addPropertiesToDatabase(notion, ids, labels)
	addCardToRespectiveDatabase(config, notion, ids, cards)
}

func LoadCardsFromExport(exportPath string) []Card {
	log.Printf("Loading cards from export path %s\n", exportPath)

	jsonFile, err := os.Open(exportPath)
	if err != nil {
		log.Fatalf("Error occurred: %v\n", err)
	}

	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	var cards []Card
	json.Unmarshal(data, &cards)

	return cards
}

// For every card, we will add the respective properties and then add the access link attachments
func addCardToRespectiveDatabase(config *config.Config, notion *notionapi.Client, ids map[string]string, cards []Card) {
	log.Printf("Adding cards to database")

	for _, card := range cards {
		log.Printf("Adding card %s\n", card.Name)

		fileAttachments, urlAttachments := organizeAttachments(card)

		_, primaryLink := FirstItem(urlAttachments)

		properties := createProperties(card, primaryLink)
		children := createChildren(fileAttachments, urlAttachments, card.Comments)

		_, err := notion.Page.Create(context.Background(), &notionapi.PageCreateRequest{
			Parent: notionapi.Parent{
				DatabaseID: notionapi.DatabaseID(ids[config.Database[card.ParentListName]]),
			},
			Properties: properties,
			Children:   children,
		})

		if err != nil {
			log.Fatalf("Error occurred when adding cards to database: %v\n", err)
		}
	}
}

// TODO: Talk about how this type needs to be filled in explicitly as the marshalled JSON will be wrong otherwise
// TODO: Figure out how to stop this from adding in alphabetical order - can make 2 different requests
func addPropertiesToDatabase(notion *notionapi.Client, ids map[string]string, labels []Label) {
	log.Printf("Adding properties to database")

	var labelOptions []notionapi.Option
	for _, label := range labels {
		labelOptions = append(labelOptions, notionapi.Option{
			Name:  label.Name,
			Color: notionapi.Color(label.Color),
		})
	}

	for _, id := range ids {
		request := &notionapi.DatabaseUpdateRequest{
			Properties: notionapi.PropertyConfigs{
				"Description": notionapi.RichTextPropertyConfig{
					Type: notionapi.PropertyConfigTypeRichText,
				},
				"Primary Link": notionapi.URLPropertyConfig{
					Type: notionapi.PropertyConfigTypeURL,
				},
				"Labels": notionapi.MultiSelectPropertyConfig{
					Type: notionapi.PropertyConfigTypeMultiSelect,
					MultiSelect: notionapi.Select{
						Options: labelOptions,
					},
				},
				"Last Updated": notionapi.DatePropertyConfig{
					Type: notionapi.PropertyConfigTypeDate,
				},
			},
		}
		_, err := notion.Database.Update(context.Background(), notionapi.DatabaseID(id), request)

		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}
}

func createProperties(card Card, primaryLink string) notionapi.Properties {
	labelOptions := organizeLabels(card.Labels)

	properties := notionapi.Properties{}

	properties["Name"] = notionapi.TitleProperty{
		Type: notionapi.PropertyTypeTitle,
		Title: []notionapi.RichText{
			{
				Text: notionapi.Text{
					Content: card.Name,
				},
			},
		},
	}

	properties["Description"] = notionapi.RichTextProperty{
		Type: notionapi.PropertyTypeRichText,
		RichText: []notionapi.RichText{
			{
				Text: notionapi.Text{
					Content: card.Description,
				},
			},
		},
	}

	if primaryLink != "null" {
		properties["Primary Link"] = notionapi.URLProperty{
			Type: notionapi.PropertyTypeURL,
			URL:  primaryLink,
		}
	}

	if card.LastUpdate != nil {
		properties["Last Updated"] = notionapi.DateProperty{
			Type: notionapi.PropertyTypeDate,
			Date: notionapi.DateObject{
				Start: (*notionapi.Date)(card.LastUpdate),
			},
		}
	}

	if len(labelOptions) > 0 {
		properties["Labels"] = notionapi.MultiSelectProperty{
			Type:        notionapi.PropertyTypeMultiSelect,
			MultiSelect: labelOptions,
		}
	}

	return properties
}

func createChildren(fileAttachments, urlAttachments map[string]string, comments []string) []notionapi.Block {
	var children []notionapi.Block
	children = append(children, createHeading1("File Attachments"))
	children = append(children, createLinkBlocks(fileAttachments)...)
	children = append(children, createHeading1("URL Attachments"))
	children = append(children, createLinkBlocks(urlAttachments)...)
	children = append(children, createHeading1("Comments"))

	for _, comment := range comments {
		children = append(children, createParagraph(comment))
	}

	return children
}

func createParagraph(text string) notionapi.ParagraphBlock {
	return notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			Text: []notionapi.RichText{
				{
					Text: notionapi.Text{
						Content: text,
					},
				},
			},
			Children: []notionapi.Block{},
		},
	}
}

func createHeading1(title string) notionapi.Heading1Block {
	return notionapi.Heading1Block{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeHeading1,
		},
		Heading1: notionapi.Heading{
			Text: []notionapi.RichText{
				{
					Text: notionapi.Text{
						Content: title,
					},
				},
			},
		},
	}
}

func createLinkBlocks(links map[string]string) []notionapi.Block {
	var blocks []notionapi.Block
	for name, link := range links {
		block := notionapi.BulletedListItemBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeBulletedListItem,
			},
			BulletedListItem: notionapi.ListItem{
				Text: []notionapi.RichText{
					{
						Text: notionapi.Text{
							Content: name,
							Link: &notionapi.Link{
								Url: link,
							},
						},
					},
				},
				Children: []notionapi.Block{},
			},
		}

		blocks = append(blocks, block)
	}

	return blocks
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

func organizeAttachments(card Card) (fileAttachments map[string]string, urlAttachments map[string]string) {
	fileAttachments, urlAttachments = make(map[string]string), make(map[string]string)

	for _, attachment := range card.Attachments {
		if attachment.IsUpload {
			fileAttachments[attachment.Name] = attachment.Url
		} else {
			urlAttachments[attachment.Name] = attachment.Url
		}
	}

	return
}

func extractLabels(config *config.Config, cards []Card) []Label {
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

	var labels []Label
	for name, color := range labelsMap {
		labels = append(labels, Label{Name: name, Color: color})
	}

	return labels
}

func getDatabaseIdsFromNames(notion *notionapi.Client, names []string) map[string]string {
	ids := make(map[string]string)

	searchResp, _ := notion.Search.Do(context.Background(), &notionapi.SearchRequest{
		Filter: map[string]string{
			"value":    "database",
			"property": "object",
		},
	})

	for _, result := range searchResp.Results {
		r := result.(*notionapi.Database)
		title := r.Title[0].Text.Content
		if Contains(title, names) {
			ids[title] = r.ID.String()
		}
	}

	for _, name := range names {
		if _, ok := ids[name]; !ok {
			log.Fatalf("Unable to find database titled %s to import to\n", name)
		}
	}

	return ids
}

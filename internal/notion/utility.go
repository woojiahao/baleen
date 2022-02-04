package notion

import (
	"time"

	na "github.com/jomei/notionapi"
)

const noLink = "null"

func richTextConfig() na.RichTextPropertyConfig {
	return na.RichTextPropertyConfig{
		Type: na.PropertyConfigTypeRichText,
	}
}

func urlConfig() na.URLPropertyConfig {
	return na.URLPropertyConfig{
		Type: na.PropertyConfigTypeURL,
	}
}

func dateConfig() na.DatePropertyConfig {
	return na.DatePropertyConfig{
		Type: na.PropertyConfigTypeDate,
	}
}

func multiSelectConfig(options []na.Option) na.MultiSelectPropertyConfig {
	return na.MultiSelectPropertyConfig{
		Type: na.PropertyConfigTypeMultiSelect,
		MultiSelect: na.Select{
			Options: options,
		},
	}
}

func titleProperty(title string) na.TitleProperty {
	return na.TitleProperty{
		Type:  na.PropertyTypeTitle,
		Title: richText(title, noLink),
	}
}

func richTextProperty(content, link string) na.RichTextProperty {
	return na.RichTextProperty{
		Type:     na.PropertyTypeRichText,
		RichText: richText(content, link),
	}
}

func urlProperty(url string) na.URLProperty {
	return na.URLProperty{
		Type: na.PropertyTypeURL,
		URL:  url,
	}
}

func dateProperty(date *time.Time) na.DateProperty {
	return na.DateProperty{
		Type: na.PropertyTypeDate,
		Date: na.DateObject{
			Start: (*na.Date)(date),
		},
	}
}

func multiSelectProperty(options []na.Option) na.MultiSelectProperty {
	return na.MultiSelectProperty{
		Type:        na.PropertyTypeMultiSelect,
		MultiSelect: options,
	}
}

func paragraph(text string) na.ParagraphBlock {
	return na.ParagraphBlock{
		BasicBlock: na.BasicBlock{
			Object: na.ObjectTypeBlock,
			Type:   na.BlockTypeParagraph,
		},
		Paragraph: na.Paragraph{
			Text:     richText(text, noLink),
			Children: []na.Block{},
		},
	}
}

func heading1(title string) na.Heading1Block {
	return na.Heading1Block{
		BasicBlock: basicBlock(na.BlockTypeHeading1),
		Heading1: na.Heading{
			Text: richText(title, noLink),
		},
	}
}

func linkBlocks(items map[string]string) []na.BulletedListItemBlock {
	var blocks []na.BulletedListItemBlock

	for name, link := range items {
		block := na.BulletedListItemBlock{
			BasicBlock: basicBlock(na.BlockTypeBulletedListItem),
			BulletedListItem: na.ListItem{
				Text:     richText(string(name), string(link)),
				Children: []na.Block{},
			},
		}

		blocks = append(blocks, block)
	}

	return blocks
}

func toBlocks(blocks []na.Block) []na.Block {
	var b []na.Block

	for _, block := range blocks {
		b = append(b, block)
	}

	return b
}

func basicBlock(t na.BlockType) na.BasicBlock {
	return na.BasicBlock{
		Object: na.ObjectTypeBlock,
		Type:   t,
	}
}

func richText(content, link string) []na.RichText {
	if link == noLink {
		return []na.RichText{
			{
				Text: na.Text{
					Content: content,
					Link: &na.Link{
						Url: link,
					},
				},
			},
		}
	} else {
		return []na.RichText{
			{
				Text: na.Text{
					Content: content,
				},
			},
		}
	}
}

func contains(s string, arr []string) bool {
	for _, x := range arr {
		if x == s {
			return true
		}
	}

	return false
}

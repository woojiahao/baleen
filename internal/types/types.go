package types

import "time"

// TODO: Change the attachment configuration
// Special cards are cards with attachments and comments
type Card struct {
	Id             string
	Name           string
	Description    string
	ParentListName string
	Labels         []*Label
	LastUpdate     *time.Time
	IsSpecial      bool
	Comments       []string
	Attachments    []*Attachment
}

type Label struct {
	Name  string
	Color string
}

type Attachment struct {
	IsUpload bool
	Name     string
	Url      string
}

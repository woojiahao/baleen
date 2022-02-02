package notion

type NotionPropertyType string

const (
	RichText    NotionPropertyType = "rich_text"
	MultiSelect NotionPropertyType = "multi_select"
	Date        NotionPropertyType = "date"
	Url         NotionPropertyType = "url"
)

type NotionObjectType string

const (
	Paragraph NotionObjectType = "paragraph"
	Heading1  NotionObjectType = "heading_1"
	Heading2  NotionObjectType = "heading_2"
	Bookmark  NotionObjectType = "bookmark"
)

type NotionMultiSelect struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type NotionPropertyBody struct {
	MultiSelectOptions []NotionMultiSelect `json:"options,omitempty"`
}

type NotionProperty map[string]NotionPropertyBody
type NotionProperties map[string]NotionProperty

type NotionPageChildBody struct {
}
type NotionPageChild map[NotionObjectType]NotionPageChildBody

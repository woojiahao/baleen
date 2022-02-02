package notion

type NotionPropertyType string

const (
	RichText    NotionPropertyType = "rich_text"
	MultiSelect NotionPropertyType = "multi_select"
	Date        NotionPropertyType = "date"
	Url         NotionPropertyType = "url"
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

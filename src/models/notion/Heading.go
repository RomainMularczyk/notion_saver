package notionModels

type NotionHeading struct {
	RichText     []NotionRichText `json:"rich_text"`
	IsToggleable bool             `json:"is_toggleable"`
	Color        string           `json:"color"`
}

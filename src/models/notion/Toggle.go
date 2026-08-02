package notionModels

type NotionToggle struct {
	RichText []NotionRichText `json:"rich_text"`
	Color    string           `json:"color"`
	Children []NotionBlock    `json:"children"`
}

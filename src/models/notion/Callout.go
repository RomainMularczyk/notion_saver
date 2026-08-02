package notionModels

type NotionCallout struct {
	RichText []NotionRichText `json:"rich_text"`
	Icon     NotionIcon       `json:"icon"`
	Color    string           `json:"color"`
}

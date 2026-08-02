package notionModels

type NotionBlockParagraph struct {
	RichText []NotionRichText `json:"rich_text"`
	Icon     *string          `json:"icon"`
	Color    string           `json:"color"`
	Children []NotionBlock    `json:"children"`
}

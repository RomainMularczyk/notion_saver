package notionModels

type NotionQuote struct {
	RichText []NotionRichText `json:"rich_text"`
	Color    string           `json:"color"`
	Children []NotionBlock    `json:"children"`
}

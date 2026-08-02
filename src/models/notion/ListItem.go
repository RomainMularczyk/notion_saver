package notionModels

type NotionBulletedListItem struct {
	RichText []NotionRichText `json:"rich_text"`
	Color    string           `json:"color"`
	Children []NotionBlock    `json:"children"`
}

type NotionNumberedListItem struct {
	RichText       []NotionRichText `json:"rich_text"`
	Color          string           `json:"color"`
	ListStartIndex int              `json:"list_start_index"`
	ListFormat     string           `json:"list_format"`
	Children       []NotionBlock    `json:"children"`
}

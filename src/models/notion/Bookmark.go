package notionModels

type NotionBookmark struct {
	Caption []NotionRichText `json:"caption"`
	Url     string           `json:"url"`
}

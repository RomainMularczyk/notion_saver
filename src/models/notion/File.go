package notionModels

type NotionFile struct {
	Caption  []NotionRichText       `json:"caption"`
	Type     string                 `json:"type"`
	External NotionFileExternalLink `json:"external"`
	Name     string                 `json:"name"`
}

type NotionFileExternalLink struct {
	Url string `json:"url"`
}

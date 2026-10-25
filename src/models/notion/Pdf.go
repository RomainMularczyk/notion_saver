package notionModels

type NotionPdf struct {
	Name     string                `json:"name"`
	External NotionPdfExternalLink `json:"external"`
}

type NotionPdfExternalLink struct {
	Url string `json:"url"`
}

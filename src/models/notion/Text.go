package notionModels

type NotionText struct {
	Content string `json:"content"`
	Link    *struct {
		Url string `json:"url"`
	}
}

type NotionRichText struct {
	Type        string           `json:"type"`
	Text        NotionText       `json:"text"`
	Annotations NotionAnnotation `json:"annotations"`
	PlainText   string           `json:"plain_text"`
	Href        string           `json:"href"`
}

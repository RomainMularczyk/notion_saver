package notionModels

type NotionImage struct {
	Caption  []NotionImageCaption `json:"caption"`
	Type     string               `json:"type"`
	External struct {
		Url string `json:"url"`
	}
}

type NotionImageCaption struct {
	Type        string           `json:"type"`
	Text        NotionText       `json:"text"`
	Annotations NotionAnnotation `json:"annotations"`
	PlainText   string           `json:"plain_text"`
	Href        string           `json:"href"`
}

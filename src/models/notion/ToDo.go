package notionModels

type NotionToDo struct {
	RichText []NotionRichText `json:"rich_text"`
	Checkbox bool             `json:"checkbox"`
	Color    string           `json:"color"`
	Children []NotionBlock    `json:"children"`
}

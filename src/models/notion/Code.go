package notionModels

type NotionCode struct {
	Caption  []NotionCodeCaption `json:"caption"`
	RichText []NotionRichText    `json:"rich_text"`
	Language string              `json:"language"`
}

type NotionCodeCaption struct {
}

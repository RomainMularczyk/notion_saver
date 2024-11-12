package models

import (
	"github.com/google/uuid"
	"time"
)

// NotionSaver models
type Blocks struct {
	Id       uuid.UUID   `gorm:"primaryKey" json:"id"`
	FullText string      `json:"full_text"`
	PageId   uuid.UUID   `json:"page_id"`
	Data     []BlockData `gorm:"foreignKey:BlocksId;constraint:OnDelete:CASCADE;" json:"data"`
}

type BlockData struct {
	Id             uuid.UUID    `gorm:"primaryKey" json:"id"`
	Object         string       `json:"object"`
	LastEditedTime time.Time    `json:"last_edited_time"`
	Type           string       `json:"type"`
	PlainText      string       `json:"plain_text"`
	BlocksId       uuid.UUID    `json:"blocks_id"`
	Annotation     []Annotation `gorm:"foreignKey:BlockId;constraint:OnDelete:CASCADE;" json:"annotations"`
}

// Notion API models
type NotionBlocks struct {
	Results []NotionBlock `json:"results"`
}

type NotionBlock struct {
	Object         string               `json:"object"`
	Id             string               `json:"id"`
	LastEditedTime time.Time            `json:"last_edited_time"`
	HasChildren    bool                 `json:"has_children"`
	Type           string               `json:"type"`
	Paragraph      NotionBlockParagraph `json:"paragraph"`
}

type NotionBlockParagraph struct {
	RichText []NotionBlockParagraphText `json:"rich_text"`
}

type NotionBlockParagraphText struct {
	Type        string                             `json:"type"`
	Annotations NotionBlockParagraphTextAnnotation `json:"annotations"`
	PlainText   string                             `json:"plain_text"`
}

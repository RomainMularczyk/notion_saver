package models

import (
	"github.com/google/uuid"
)

// NotionSaver models
type Annotation struct {
	Id            uuid.UUID `gorm:"primaryKey" json:"id"`
	Text          string    `json:"text"`
	Bold          bool      `json:"bold"`
	Italic        bool      `json:"italic"`
	Strikethrough bool      `json:"strikethrough"`
	Underline     bool      `json:"underline"`
	Code          bool      `json:"code"`
	Color         string    `json:"color"`
	BlockId       uuid.UUID `json:"block_id"`
}

// Notion models
type NotionBlockParagraphTextAnnotation struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

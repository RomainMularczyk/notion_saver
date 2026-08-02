package saverModels

import (
	"github.com/google/uuid"
	"time"
)

// NotionSaver models
type Page struct {
	Id         uuid.UUID `gorm:"primaryKey" validate:"required" json:"id"`
	Title      string    `gorm:"not null,index:idx_page_title" validate:"required" json:"title"`
	PageType   string    `gorm:"not null" validate:"required" json:"page_type"`
	LastEdited time.Time `gorm:"not null" validate:"required" json:"last_edited"`
	EmojiIcon  string    `json:"emoji_icon"`
	IconLink   string    `json:"icon_link"`
	Saves      []Save    `gorm:"many2many:save_pages;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Blocks     []Block   `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

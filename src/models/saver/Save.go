package saverModels

import (
	"github.com/google/uuid"
	"time"
)

type SaveStatus string

const (
	Pending    SaveStatus = "pending"
	InProgress SaveStatus = "in_progress"
	Completed  SaveStatus = "completed"
	Failed     SaveStatus = "failed"
)

type Save struct {
	Id       uuid.UUID  `gorm:"primaryKey" json:"id"`
	LastSave time.Time  `json:"last_save"`
	Status   SaveStatus `json:"status"`
	Pages    []Page     `gorm:"many2many:save_pages;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

// The many-to-many relationship between a Save and a Page
type SavePage struct {
	SaveId uuid.UUID `gorm:"primaryKey;uniqueIndex:idx_save_page" json:"save_id"`
	PageId uuid.UUID `gorm:"primaryKey;uniqueIndex:idx_save_page" json:"page_id"`

	Page Page `gorm:"foreignKey:PageId;references:Id"`
	Save Save `gorm:"foreignKey:SaveId;references:Id"`
}

package models

import (
	"github.com/google/uuid"
	"time"
)

type Save struct {
	Id       uuid.UUID `gorm:"primaryKey" json:"id"`
	LastSave time.Time `json:"last_save"`
	Pages    []Page    `gorm:"many2many:save_page;"`
}

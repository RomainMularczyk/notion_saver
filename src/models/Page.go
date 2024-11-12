package models

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

// NotionSaver models
type PageType int

const (
	SimplePage PageType = iota
	Database
)

type Page struct {
	Id         uuid.UUID `gorm:"primaryKey" json:"id"`
	Title      string    `json:"title"`
	PageType   string    `json:"page_type"`
	LastEdited time.Time `json:"last_edited"`
	EmojiIcon  string    `json:"emoji_icon"`
	IconLink   string    `json:"icon_link"`
	Hash       *string   `json:"hash"`
	Blocks     []Blocks  `gorm:"foreignKey:PageId;constraint:OnDelete:CASCADE;" json:"blocks"`
	SaveId     uuid.UUID `json:"save_id"`
}

// Notion API models
type NotionParentPage struct {
	PageType PageType `json:"type"`
	PageId   string   `json:"page_id"`
}

type NotionTitle struct {
	Id    uuid.UUID               `json:"id"`
	Type  string                  `json:"type"`
	Title []NotionTitleProperties `json:"title"`
}

type NotionTitleProperties struct {
	PlainText string `json:"plain_text"`
}

type NotionPageProperties struct {
	Title NotionTitle `json:"title"`
}

type NotionPageIcon struct {
	Type     string                     `json:"type"`
	External NotionPageIconExternalLink `json:"external"`
	Emoji    string                     `json:"emoji"`
	File     NotionPageIconFile         `json:"file"`
}

type NotionPageIconExternalLink struct {
	Url string `json:"url"`
}

type NotionPageIconFile struct {
	Url string `json:"url"`
}

type NotionPage struct {
	Object         string               `json:"object"`
	Id             uuid.UUID            `json:"id"`
	LastEditedTime time.Time            `json:"last_edited_time"`
	Parent         NotionParentPage     `json:"parent"`
	Icon           NotionPageIcon       `json:"icon"`
	Properties     NotionPageProperties `json:"properties"`
}

type NotionPages struct {
	Results []NotionPage `json:"results"`
}

// Convert a Notion API Page into a NotionSaver Page
func ToPage(save_id string, data NotionPages) []Page {
	var pages []Page

	for _, result := range data.Results {
		// Check for pages with empty titles
		if len(result.Properties.Title.Title) == 0 {
			slog.Warn(fmt.Sprintf("Page did not have title. Skipping page with Id: %v", result.Id))
		} else {
			if result.Icon.Emoji != "" {
				page := Page{
					Id:         result.Id,
					Title:      result.Properties.Title.Title[0].PlainText,
					PageType:   result.Object,
					LastEdited: result.LastEditedTime,
					EmojiIcon:  result.Icon.Emoji,
				}
				pages = append(pages, page)
			} else {
				if result.Icon.Type == "external" {
					page := Page{
						Id:         result.Id,
						Title:      result.Properties.Title.Title[0].PlainText,
						PageType:   result.Object,
						LastEdited: result.LastEditedTime,
						IconLink:   result.Icon.External.Url,
					}
					pages = append(pages, page)
				}
				if result.Icon.Type == "file" {
					page := Page{
						Id:         result.Id,
						Title:      result.Properties.Title.Title[0].PlainText,
						PageType:   result.Object,
						LastEdited: result.LastEditedTime,
						IconLink:   result.Icon.File.Url,
					}
					pages = append(pages, page)
				}
			}
		}
	}

	return pages
}

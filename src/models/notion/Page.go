package notionModels

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"

	saverModels "notion_saver/src/models/saver"
)

// Notion API models
type NotionParentPage struct {
	PageType string `json:"type"`
	PageId   string `json:"page_id"`
}

type NotionTitle struct {
	Id    string                  `json:"id"`
	Type  string                  `json:"type"`
	Title []NotionTitleProperties `json:"title"`
}

type NotionTitleProperties struct {
	Type        string           `json:"type"`
	Annotations NotionAnnotation `json:"annotations"`
	PlainText   string           `json:"plain_text"`
	Href        string           `json:"href"`
}

type NotionPageProperties struct {
	Title NotionTitle `json:"title"`
}

type NotionPageIcon struct {
	Type        string                     `json:"type"`
	External    NotionPageIconExternalLink `json:"external"`
	Emoji       string                     `json:"emoji"`
	CustomEmoji NotionCustomEmoji          `json:"custom_emoji"`
	File        NotionPageIconFile         `json:"file"`
}

type NotionCustomEmoji struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type NotionPageIconExternalLink struct {
	Url string `json:"url"`
}

type NotionPageIconFile struct {
	Url string `json:"url"`
}

type NotionPage struct {
	Object         string                  `validate:"required" json:"object"`
	Id             uuid.UUID               `validate:"required" json:"id"`
	CreatedTime    time.Time               `validate:"required" json:"created_time"`
	LastEditedTime time.Time               `validate:"required" json:"last_edited_time"`
	DatabaseTitle  []NotionTitleProperties `json:"title"`
	Icon           NotionPageIcon          `validate:"required" json:"icon"`
	Parent         NotionParentPage        `json:"parent"`
	InTrash        bool                    `json:"in_trash"`
	IsArchived     bool                    `json:"is_archived"`
	IsLocked       bool                    `json:"is_locked"`
	Properties     NotionPageProperties    `validate:"required" json:"properties"`
	Url            string                  `json:"url"`
	PublicUrl      string                  `json:"public_url"`
}

type NotionPages struct {
	Object     string       `json:"object"`
	Results    []NotionPage `json:"results"`
	NextCursor string       `json:"next_cursor"`
	HasMore    bool         `json:"has_more"`
	Type       string       `json:"type"`
}

// Convert a Notion API Page into a NotionSaver Page
func ToPage(save_id string, data NotionPages) []saverModels.Page {
	var pages []saverModels.Page

	for _, result := range data.Results {
		// Check for pages with empty titles
		if len(result.Properties.Title.Title) == 0 {
			slog.Warn(fmt.Sprintf("Page did not have title. Skipping page with Id: %v", result.Id))
		} else {
			if result.Icon.Emoji != "" {
				page := saverModels.Page{
					Id:         result.Id,
					Title:      result.Properties.Title.Title[0].PlainText,
					PageType:   result.Object,
					LastEdited: result.LastEditedTime,
					EmojiIcon:  result.Icon.Emoji,
				}
				pages = append(pages, page)
			} else {
				if result.Icon.Type == "external" {
					page := saverModels.Page{
						Id:         result.Id,
						Title:      result.Properties.Title.Title[0].PlainText,
						PageType:   result.Object,
						LastEdited: result.LastEditedTime,
						IconLink:   result.Icon.External.Url,
					}
					pages = append(pages, page)
				}
				if result.Icon.Type == "file" {
					page := saverModels.Page{
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

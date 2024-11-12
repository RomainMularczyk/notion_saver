package data

import (
	"fmt"
	"github.com/google/uuid"
	"notion_saver/src/models"
	"strings"
)

// Gather a Notion page full text content
func BuildPageContent(pageContent models.NotionBlocks) string {
	var textContent strings.Builder
	for _, result := range pageContent.Results {
		if result.Type == "paragraph" {
			textContent.WriteString(
				BuildPageContentFromBlock(result),
			)
		}
	}
	return textContent.String()
}

// Retrieve page metadata from Notion page
func BuildPageFromNotionPage(
	notionPage models.NotionPage,
	saveId uuid.UUID,
) models.Page {
	if notionPage.Object == "page" {
		if notionPage.Icon.Type == "external" {
			page := models.Page{
				Id:         notionPage.Id,
				Title:      notionPage.Properties.Title.Title[0].PlainText,
				PageType:   notionPage.Properties.Title.Type,
				LastEdited: notionPage.LastEditedTime,
				EmojiIcon:  notionPage.Icon.Emoji,
				IconLink:   notionPage.Icon.External.Url,
				SaveId:     saveId,
			}
			return page
		} else if notionPage.Icon.Type == "file" {
			page := models.Page{
				Id:         notionPage.Id,
				Title:      notionPage.Properties.Title.Title[0].PlainText,
				PageType:   notionPage.Properties.Title.Type,
				LastEdited: notionPage.LastEditedTime,
				EmojiIcon:  notionPage.Icon.Emoji,
				IconLink:   notionPage.Icon.File.Url,
				SaveId:     saveId,
			}
			return page
		} else {
			page := models.Page{
				Id:         notionPage.Id,
				Title:      notionPage.Properties.Title.Title[0].PlainText,
				PageType:   notionPage.Properties.Title.Type,
				LastEdited: notionPage.LastEditedTime,
				EmojiIcon:  notionPage.Icon.Emoji,
				SaveId:     saveId,
			}
			return page
		}
	} else {
		return BuildDatabaseFromNotionPage(notionPage)
	}
}

func BuildDatabaseFromNotionPage(
	notionPage models.NotionPage,
	saveId uuid.UUID,
) models.Page {
}

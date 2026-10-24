package utils

import (
	"log/slog"
	notionModels "notion_saver/src/models/notion"
	saverModels "notion_saver/src/models/saver"
)

// Create SaverPages from the NotionPages retrieved from the Notion API
func BuildSaverPages(notionPages *notionModels.NotionPages) []saverModels.Page {
	newPages := make([]saverModels.Page, 0, len(notionPages.Results))

	for _, page := range notionPages.Results {
		if len(page.Properties.Title.Title) == 0 {
			// slog.Warn(
			// 	fmt.Sprintf(
			// 		"Page did not have title. Skipping page with Id: %v", page.Id,
			// 	),
			// )
			continue
		}

		newPage := saverModels.Page{
			Id:         page.Id,
			Title:      page.Properties.Title.Title[0].PlainText,
			PageType:   page.Object,
			LastEdited: page.LastEditedTime,
			EmojiIcon:  page.Icon.Emoji,
			IconLink:   page.Icon.External.Url,
		}
		newPages = append(newPages, newPage)
		slog.Debug(
			"New page formed successfully.",
			slog.String("Page", newPage.Title),
		)
	}

	slog.Debug("Pages were built successfully.", slog.Int("pages_counter", len(newPages)))

	return newPages
}

package services

import (
	"notion_saver/src/models"
	"notion_saver/src/queries"
	"notion_saver/src/repositories"
	"notion_saver/src/utils/data"
	"notion_saver/src/utils/hashing"
)

func CreatePage(
	page models.Page,
) (*models.Page, error) {
	// Create page
	queryPage := repositories.QueryPage()
	result, err := queryPage.Create(page)
	if err != nil {
		return nil, err
	}
	// Retrieve Notion blocks contained on the page
	notionBlocks := queries.GetNotionBlocks(page.Id)
	pageContentAsText := data.BuildPageContent(notionBlocks)
	// Create a hash representing the text content of the page
	pageHash := hashing.HashPage(pageContentAsText)
	page.Hash = &pageHash
	// Create the blocks contained in the page
	CreateBlocks(notionBlocks, page.Id)
	return result, nil
}

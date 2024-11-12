package services

import (
	"notion_saver/src/models"
	"notion_saver/src/queries"
	"notion_saver/src/repositories"
	"notion_saver/src/utils/data"
)

func CreateSave(
	save models.Save,
) (*models.Save, error) {
	// Create save
	querySave := repositories.QuerySave()
	result, err := querySave.Create(save)
	// Retrieve Notion pages
	notionPages := queries.GetAllPages()
	for _, notionPage := range notionPages.Results {
		page := data.BuildPageFromNotionPage(notionPage, save.Id)
		CreatePage(page)
	}
	return result, err
}

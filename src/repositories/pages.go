package repositories

import (
	"fmt"
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/utils"

	"github.com/google/uuid"
)

type PageRepository struct {
	server *utils.Server
	logger *slog.Logger
}

func NewPageRepository(server *utils.Server) *PageRepository {
	return &PageRepository{
		server: server,
		logger: server.Logger.With("module", "repositories/pages"),
	}
}

// Save many pages in the database
// Also creates the relation between the pages and the save
func (p *PageRepository) CreateMany(
	server *utils.Server,
	pages []saverModels.Page,
	saveId uuid.UUID,
) (*[]saverModels.Page, error) {
	p.logger.Debug(
		fmt.Sprintf("Attempting to create %v pages.", len(pages)),
	)
	tx := server.Database.Begin()

	if tx.Error != nil {
		p.logger.Error(fmt.Sprintf("Error when starting a transaction: %v", tx.Error))
		return nil, tx.Error
	}

	var savePages []saverModels.SavePage
	for _, page := range pages {
		savePage := saverModels.SavePage{
			PageId: page.Id,
			SaveId: saveId,
		}
		savePages = append(savePages, savePage)
	}

	if err := tx.Create(&pages).Error; err != nil {
		tx.Rollback()
		p.logger.Error(fmt.Sprintf("Error when creating pages: %v", err))
		return nil, err
	}
	if err := tx.Create(&savePages).Error; err != nil {
		tx.Rollback()
		p.logger.Error(fmt.Sprintf("Error when creating save pages: %v", err))
		return nil, err
	}
	tx.Commit()

	p.logger.Debug(fmt.Sprintf("%v pages created successfully.", len(pages)))

	return &pages, nil
}

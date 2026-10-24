package services

import (
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/repositories"
	"notion_saver/src/utils"

	"github.com/google/uuid"
)

type PageService struct {
	server *utils.Server
	logger *slog.Logger
}

func NewPageService(server *utils.Server) *PageService {
	return &PageService{
		server: server,
		logger: server.Logger.With("module", "services/pages"),
	}
}

// Save many pages in the database
func (p *PageService) SaveMany(
	server *utils.Server,
	pages []saverModels.Page,
	saveId uuid.UUID,
) (*[]saverModels.Page, error) {
	pageRepository := repositories.NewPageRepository(server)
	savedPages, err := pageRepository.CreateMany(server, pages, saveId)
	if err != nil {
		return nil, err
	}
	p.logger.Debug("Pages saved successfully.", slog.Int("Number of pages", len(*savedPages)))
	return savedPages, nil
}

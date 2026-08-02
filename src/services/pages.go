package services

import (
	"fmt"
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
	p.logger.Debug(fmt.Sprintf("%v pages saved successfully.", len(*savedPages)))
	return savedPages, nil
}

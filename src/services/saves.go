package services

import (
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/repositories"
	"notion_saver/src/utils"
	"time"

	"github.com/google/uuid"
)

type SaveService struct {
	server *utils.Server
	logger *slog.Logger
}

func NewSaveService(server *utils.Server) *SaveService {
	return &SaveService{
		server: server,
		logger: server.Logger.With("module", "services/saves"),
	}
}

// Create a new save. This is the main entry point of the save mechanism.
// Starting a new save will:
// 1. Create a new save object in the database
// 2. Debezium will emit a message to the RabbitMQ queue so the worker
// can start crawling the Notion API
func (s *SaveService) New(server *utils.Server) (*saverModels.Save, error) {
	newSaveId := uuid.New()
	save := saverModels.Save{
		Id:       newSaveId,
		LastSave: time.Now(),
	}

	saveRepository := repositories.NewSaveRepository(server)
	newSave, err := saveRepository.Create(server, save)
	if err != nil {
		return nil, err
	}

	return newSave, nil
}

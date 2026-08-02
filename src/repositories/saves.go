package repositories

import (
	"fmt"
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/utils"
)

type SaveRepository struct {
	server *utils.Server
	logger *slog.Logger
}

func NewSaveRepository(server *utils.Server) *SaveRepository {
	return &SaveRepository{
		server: server,
		logger: server.Logger.With("module", "repositories/saves"),
	}
}

// Create a new save in the database
func (s *SaveRepository) Create(
	server *utils.Server,
	save saverModels.Save,
) (*saverModels.Save, error) {
	s.logger.Debug(
		fmt.Sprintf("Attempting to create a new save."),
	)
	result := server.Database.Create(&save)
	if result.Error != nil {
		s.logger.Error(fmt.Sprintf("Error when creating save: %v", result.Error))
		return nil, result.Error
	}
	s.logger.Debug(fmt.Sprintf("Save created successfully."))

	return &save, nil
}

package repositories

import (
	"fmt"
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/utils"
)

type BlockRepository struct {
	server *utils.Server
	logger *slog.Logger
}

func NewBlockRepository(server *utils.Server) *BlockRepository {
	return &BlockRepository{
		server: server,
		logger: server.Logger.With("module", "repositories/blocks"),
	}
}

// Create many blocks in the database
func (b *BlockRepository) CreateMany(
	server *utils.Server,
	blocks []saverModels.Block,
) (*[]saverModels.Block, error) {
	b.logger.Debug(
		fmt.Sprintf("Attempting to create %v blocks.", len(blocks)),
	)

	result := server.Database.Create(&blocks)
	if err := result.Error; err != nil {
		b.logger.Error(fmt.Sprintf("Error when creating blocks: %v", err))
		return nil, err
	}
	b.logger.Debug(fmt.Sprintf("%v blocks created successfully.", len(blocks)))

	return &blocks, nil
}

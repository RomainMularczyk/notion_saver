package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
	"notion_saver/src/database"
	"notion_saver/src/models"
)

type BlockRepository struct {
	DB *gorm.DB
}

func QueryBlock() *BlockRepository {
	db := database.DatabaseConfig()
	return &BlockRepository{DB: db}
}

// Retrieve one blocks in the database
func (br *BlockRepository) RetrieveOne(id string) models.Blocks {
	var blocks models.Blocks

	uuid, err := uuid.Parse(id)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Failed to parse UUID: %v", err),
		)
	}

	result := br.DB.Model(models.Blocks{Id: uuid}).First(&blocks)

	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf(
				"An error occurred when trying to retrieve Blocks with Id: %s",
				id,
			),
		)
	}
	return blocks
}

// Insert a new blocks in the database
func (br *BlockRepository) Create(blocks models.Blocks) (*models.Blocks, error) {
	result := br.DB.Create(&blocks)
	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to create a new blocks: %s", err),
		)
		return nil, err
	}

	return &blocks, nil
}

// Insert a new block in the database
func (br *BlockRepository) CreateBlockData(block models.BlockData) (*models.BlockData, error) {
	result := br.DB.Create(&block)
	if err := result.Error; err != nil {
		slog.Error(
			fmt.Sprintf("An error occurred when trying to create a new block: %s", err),
		)
		return nil, err
	} else {
		slog.Info(
			fmt.Sprintf(
				"New block was inserted successfully: %s", block.Id,
			),
		)
		return &block, nil
	}
}

// Insert many blocks in the database
func (br *BlockRepository) CreateBlocksData(blocks []models.BlockData) (*[]models.BlockData, error) {
	var blocksData []models.BlockData
	for _, block := range blocks {
		block, err := br.CreateBlockData(block)
		if err != nil {
			return nil, err
		}
		blocksData = append(blocksData, *block)
	}
	return &blocksData, nil
}

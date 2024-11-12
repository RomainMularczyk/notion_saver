package services

import (
	"github.com/google/uuid"
	"notion_saver/src/models"
	"notion_saver/src/repositories"
	"notion_saver/src/utils/data"
)

func CreateBlocks(
	notionBlocks models.NotionBlocks,
	pageId uuid.UUID,
) (*models.Blocks, error) {
	// Create blocks
	blocks := data.BuildBlocksContent(notionBlocks, pageId)
	queryBlock := repositories.QueryBlock()
	result, err := queryBlock.Create(blocks)
	// Create block data
	CreateBlockData(notionBlocks, result.Id)
	return result, err
}

func CreateBlockData(notionBlocks models.NotionBlocks, blocksId uuid.UUID) error {
	// Create block data
	blocksData, err := data.BuildBlocksData(notionBlocks, blocksId)
	queryBlock := repositories.QueryBlock()
	queryBlock.CreateBlocksData(*blocksData)
	return err
}

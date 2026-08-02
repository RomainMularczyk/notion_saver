package services

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/repositories"
	"notion_saver/src/utils"
	blocksWorkerModels "notion_saver/src/workers/blocks/models"
)

type BlockService struct {
	server *utils.Server
	logger *slog.Logger
}

func NewBlockService(server *utils.Server) *BlockService {
	return &BlockService{
		server: server,
		logger: server.Logger.With("module", "services/blocks"),
	}
}

// Publish a message in the queue to save all the blocks of a page
func (b *BlockService) BlocksOfPageMessage(
	server *utils.Server,
	pageId uuid.UUID,
	paginationIndex int,
) error {
	ch, err := server.Queue.Channel()
	if err != nil {
		b.logger.Error(fmt.Sprintf("Error when creating a channel: %v", err))
		return err
	}
	defer ch.Close()

	message := blocksWorkerModels.BlockMessage{
		PageId:          pageId,
		PaginationIndex: paginationIndex,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		b.logger.Error(fmt.Sprintf("Error when marshalling JSON: %v", err))
		return err
	}

	ch.Publish(
		"notion.blocks",
		"notion.blocks",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
	b.logger.Debug(
		"Message published successfully.",
		slog.String("message", utils.ToJson(message)),
	)

	return nil
}

// Save many blocks in the database
func (b *BlockService) SaveMany(
	server *utils.Server,
	blocks []saverModels.Block,
) (*[]saverModels.Block, error) {
	blockRepository := repositories.NewBlockRepository(server)
	savedBlocks, err := blockRepository.CreateMany(server, blocks)
	if err != nil {
		b.logger.Error(fmt.Sprintf("Error when saving blocks: %v", err))
		return nil, err
	}

	b.logger.Debug(fmt.Sprintf("%v blocks saved successfully.", len(*savedBlocks)))
	return savedBlocks, nil
}

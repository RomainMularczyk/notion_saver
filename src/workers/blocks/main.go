package blocksWorker

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"notion_saver/src/services"
	"notion_saver/src/utils"
	blocksWorkerModels "notion_saver/src/workers/blocks/models"
)

type BlockWorker struct {
	server *utils.Server
	logger *slog.Logger
}

func NewBlockWorker(server *utils.Server) *BlockWorker {
	return &BlockWorker{
		server: server,
		logger: server.Logger.With("module", "workers/blocks"),
	}
}

func (b *BlockWorker) Run(server *utils.Server) {
	defer server.Queue.Close()

	ch, err := server.Queue.Channel()
	if err != nil {
		panic(err)
	}

	b.logger.Info("Start consuming messages...")

	messages, err := ch.Consume(
		"notion.blocks",
		"notion.blocks",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	var forever chan blocksWorkerModels.BlockMessage
	go func() {
		for d := range messages {
			var blockMessage blocksWorkerModels.BlockMessage
			err := json.Unmarshal(d.Body, &blockMessage)
			if err != nil {
				b.logger.Error(fmt.Sprintf("Error when unmarshalling block message: %v", err))
				ch.Nack(d.DeliveryTag, false, true)
				continue
			}

			// Publish the next cursor message in the queue
			newBlocks, err := b.GetNotionBlocks(ch, blockMessage)
			if err != nil {
				b.logger.Error(fmt.Sprintf("Error when retrieving blocks from Notion API: %v", err))
				ch.Nack(d.DeliveryTag, false, true)
				continue
			}

			blockService := services.NewBlockService(server)
			blockService.SaveMany(server, newBlocks)
			ch.Ack(d.DeliveryTag, false)
		}
	}()
	<-forever
}

package blocksWorker

import (
	"encoding/json"
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
				b.logger.Error(
					"Error when unmarshalling block message.",
					slog.String("error", err.Error()),
				)
				ch.Nack(d.DeliveryTag, false, false)
				continue
			}

			// Publish the next cursor message in the queue
			newBlocks, err := b.GetNotionBlocks(ch, blockMessage)
			notion := utils.NewNotion(b.server)

			if err != nil {
				notionError := notion.IsNotionError(err)
				if notionError != nil && notionError.IsRateLimited {
					b.logger.Error(
						"Error when retrieving blocks from Notion API.",
						slog.String("error", err.Error()),
						slog.Bool("rate_limited", notionError.IsRateLimited),
						slog.Bool("client_error", notionError.IsClientError),
					)
					b.SendBlockMessageToDql(ch, blockMessage)
				}
				ch.Nack(d.DeliveryTag, false, false)
				continue
			}

			blockService := services.NewBlockService(server)
			blockService.SaveMany(server, newBlocks)
			ch.Ack(d.DeliveryTag, false)
		}
	}()
	<-forever
}

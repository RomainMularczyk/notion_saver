package blocksWorker

import (
	"log/slog"
	"notion_saver/src/queue"
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

func (w *BlockWorker) Run(server *utils.Server) {
	defer server.Queue.Close()

	ch, err := server.Queue.Channel()
	if err != nil {
		panic(err)
	}

	w.logger.Info("Start consuming messages...")

	messages, err := ch.Consume(
		"notion.blocks",
		"notion.blocks",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	retryPolicy := queue.NewRetryPolicy("notion.blocks", queue.DefaultRetryDelays)

	var forever chan blocksWorkerModels.BlockMessage
	go func() {
		for d := range messages {
			w.logger.Debug(
				"Received message from queue.",
				slog.String("exchange", retryPolicy.ExchangeName()),
			)
			blockMsg, err := w.ReadNotionBlocksMessage(d)
			if err != nil {
				w.logger.Error(
					"Message is not processable. Dropping.",
					slog.String("error", err.Error()),
				)
				ch.Nack(d.DeliveryTag, false, false)
				continue
			}

			// Publish the next cursor message in the queue
			newBlocks, notionErr := w.GetNotionBlocks(ch, blockMsg)
			switch {
			case notionErr == nil:
				blockService := services.NewBlockService(server)
				blockService.SaveMany(server, newBlocks)
				ch.Ack(d.DeliveryTag, false)
			case notionErr.IsRetryable():
				err := retryPolicy.Schedule(ch, d)
				// NOTE: if publishing to delay queue fails, we send the message to DLQ
				if err != nil {
					w.logger.Info(
						"Message could not be requeued. Dropping.",
						slog.String("error", err.Error()),
					)
					d.Nack(false, false)
				} else {
					w.logger.Info(
						"Message requeued successfully.",
						slog.String("queue", retryPolicy.ExchangeName()),
					)
					d.Ack(false)
				}
			default:
				d.Nack(false, false)
			}
		}
	}()
	<-forever
}

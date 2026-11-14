package pagesWorker

import (
	"log/slog"
	"notion_saver/src/queue"
	"notion_saver/src/services"
	"notion_saver/src/utils"
	blocksWorkerModels "notion_saver/src/workers/blocks/models"
	sharedWorker "notion_saver/src/workers/shared"
)

func NewPageWorker(server *utils.Server) *PageWorker {
	return &PageWorker{
		server: server,
		logger: server.Logger.With("module", "workers/pages"),
	}
}

// The PageWorker consumes messages from the 'notion.pages' queue
// and retrieves following pages from the Notion API
func (w *PageWorker) Run(server *utils.Server) {
	defer server.Queue.Close()

	ch, err := server.Queue.Channel()
	if err != nil {
		panic(err)
	}

	w.logger.Info("Start consuming messages...")

	messages, err := ch.Consume(
		"notion.pages",
		"notion.pages",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	retryPolicy := queue.NewRetryPolicy("notion.pages", queue.DefaultRetryDelays)

	var forever chan struct{}
	go func() {
		for d := range messages {
			w.logger.Info("Received message from notion.pages queue.")
			pageMsg, err := w.ReadNotionPagesMessage(d)
			if err != nil {
				w.logger.Error("Message is not processable. Dropping.")
				d.Nack(false, false)
				continue
			}

			notionPages, notionErr := sharedWorker.GetNextPages(ch, pageMsg, w.server)
			switch {
			case notionErr == nil:
				pageService := services.NewPageService(w.server)
				pageService.SaveMany(w.server, notionPages, pageMsg.SaveId)
				w.logger.Info("Message was successfully processed.")

				for _, notionPage := range notionPages {
					blockMsg := blocksWorkerModels.BlockMessage{
						NextCursor:      nil,
						PageId:          notionPage.Id,
						PaginationIndex: pageMsg.PaginationIndex + 1,
					}
					w.PublishBlockMessage(ch, blockMsg)
				}
				d.Ack(false)
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

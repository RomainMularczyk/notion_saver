package savesWorker

import (
	"notion_saver/src/services"
	"notion_saver/src/utils"
	sharedWorker "notion_saver/src/workers/shared"
)

func NewSaveWorker(server *utils.Server) *SaveWorker {
	return &SaveWorker{
		server: server,
		logger: server.Logger.With("module", "workers/saves"),
	}
}

// The PageWorker consumes messages from the 'notion.pages' queue
// and retrieves following pages from the Notion API
func (w *SaveWorker) Run(server *utils.Server) {
	defer server.Queue.Close()

	ch, err := server.Queue.Channel()
	if err != nil {
		panic(err)
	}

	w.logger.Info("Start consuming messages...")

	messages, err := ch.Consume(
		"debezium",
		"debezium",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	var forever chan struct{}
	go func() {
		for d := range messages {
			w.logger.Debug("Received message from Debezium queue.")
			saveMsg, err := w.ReadDebeziumSaveMessage(d)
			if err != nil {
				w.logger.Error("Message is not processable. Dropping.")
				d.Ack(false)
			} else {
				pageMessage := sharedWorker.PageMessage{
					SaveId:          saveMsg.Payload.After.Id,
					PaginationIndex: 0,
					NextCursor:      nil,
					RetryCount:      0,
				}
				notionPages, notionErr := sharedWorker.GetNextPages(
					ch,
					&pageMessage,
					w.server,
				)
				if notionErr.IsRetryable() {
					w.logger.Error("Message could not be processed. Sending to DLQ.")
					d.Nack(false, false)
				} else {
					w.logger.Error("Processing message is not retriable. Dropping.")
					d.Ack(false)
				}

				saveService := services.NewPageService(w.server)
				saveService.SaveMany(
					w.server,
					notionPages,
					saveMsg.Payload.After.Id,
				)
				w.logger.Info("Message was processed successfully.")
				d.Ack(false)
			}
		}
	}()

	<-forever
}

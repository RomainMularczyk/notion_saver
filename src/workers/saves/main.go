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
			w.logger.Info("Received message from Debezium queue.")
			saveMsg, err := w.ReadDebeziumSaveMessage(d)
			if err != nil {
				w.logger.Error(
					"An inrecoverable error occurred when handling the message. Dropping the message.",
				)
				d.Nack(false, false)
			} else {
				notionPages, err := sharedWorker.GetNextPages(
					ch,
					&sharedWorker.PageMessage{
						SaveId:          saveMsg.Payload.After.Id,
						PaginationIndex: 0,
						NextCursor:      nil,
					},
					w.server,
				)
				if err != nil {
					w.logger.Error("Could not save the pages retrieved from the Notion API.")
					d.Nack(false, false)
					continue
				}
				saveService := services.NewPageService(w.server)
				saveService.SaveMany(w.server, notionPages, saveMsg.Payload.After.Id)
				w.logger.Info("Message was successfully processed")
				d.Ack(true)
			}
		}
	}()

	<-forever
}

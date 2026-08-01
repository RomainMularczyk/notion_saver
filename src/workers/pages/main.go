package pagesWorker

import (
	"notion_saver/src/services"
	"notion_saver/src/utils"
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

	var forever chan struct{}
	go func() {
		for d := range messages {
			w.logger.Info("Received message from notion.pages queue.")
			pageMsg, err := w.ReadNotionPagesMessage(d)
			if err != nil {
				w.logger.Error(
					"An inrecoverable error occurred when handling the message. Dropping the message.",
				)
				d.Nack(false, false)
			} else {
				notionPages, err := sharedWorker.GetNextPages(
					ch,
					pageMsg,
					w.server,
				)
				if err != nil {
					w.logger.Error("Could not save the pages retrieved from the Notion API.")
					d.Nack(false, false)
					continue
				}
				pageService := services.NewPageService(w.server)
				pageService.SaveMany(w.server, notionPages, pageMsg.SaveId)
				w.logger.Info("Message was successfully processed")
				sharedWorker.GetNextPages(ch, pageMsg, w.server)
			}
		}
	}()

	<-forever
}

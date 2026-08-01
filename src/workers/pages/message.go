package pagesWorker

import (
	"encoding/json"
	"notion_saver/src/utils"

	amqp "github.com/rabbitmq/amqp091-go"
	sharedWorker "notion_saver/src/workers/shared"
)

// Reads and parses a NotionPage message representing a page retrieved
// from the Notion API
func (w *PageWorker) ReadNotionPagesMessage(
	msg amqp.Delivery,
) (*sharedWorker.PageMessage, error) {
	body := utils.FromBase64(msg.Body)

	var notionPagesMessage sharedWorker.PageMessage
	if err := json.Unmarshal(body, &notionPagesMessage); err != nil {
		w.logger.Error("Error when unmarshalling JSON")
		return nil, err
	}

	w.logger.Info("Message was parsed successfully.")
	return &notionPagesMessage, nil
}

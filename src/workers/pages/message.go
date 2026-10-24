package pagesWorker

import (
	"encoding/json"
	"notion_saver/src/utils"

	blocksWorkerModels "notion_saver/src/workers/blocks/models"
	sharedWorker "notion_saver/src/workers/shared"

	amqp "github.com/rabbitmq/amqp091-go"
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

func (w *PageWorker) PublishBlockMessage(
	ch *amqp.Channel,
	blockMessage blocksWorkerModels.BlockMessage,
) {
	w.logger.Info("Publishing block message")
	messageBytes, err := json.Marshal(blockMessage)
	if err != nil {
		w.logger.Error("Error when marshalling JSON")
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
}

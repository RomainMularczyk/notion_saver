package blocksWorker

import (
	"encoding/json"
	"notion_saver/src/utils"
	blocksWorkerModels "notion_saver/src/workers/blocks/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *BlockWorker) ReadNotionBlocksMessage(msg amqp.Delivery) (
	*blocksWorkerModels.BlockMessage,
	error,
) {
	body := utils.FromBase64(msg.Body)

	var notionBlocksMessage blocksWorkerModels.BlockMessage
	if err := json.Unmarshal(body, &notionBlocksMessage); err != nil {
		w.logger.Error("Error when unmarshalling JSON.")
		return nil, err
	}

	w.logger.Info("Message was parsed succesfully.")
	return &notionBlocksMessage, nil
}

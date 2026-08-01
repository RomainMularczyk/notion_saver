package savesWorker

import (
	"encoding/json"
	"errors"
	debeziumModels "notion_saver/src/models/debezium"
	"notion_saver/src/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Reads and parses a Debeizum message representing a save creation
func (w *SaveWorker) ReadDebeziumSaveMessage(
	msg amqp.Delivery,
) (*debeziumModels.DebeziumMessage, error) {
	body := utils.FromBase64(msg.Body)

	var debeziumMsg debeziumModels.DebeziumMessage
	if err := json.Unmarshal(body, &debeziumMsg); err != nil {
		w.logger.Error("Error when unmarshalling JSON")
		return nil, err
	}

	if !debeziumMsg.IsValid() {
		w.logger.Error("Message is not valid")
		return nil, errors.New("Message is not valid")
	}

	w.logger.Info("Message was parsed successfully.")
	return &debeziumMsg, nil
}

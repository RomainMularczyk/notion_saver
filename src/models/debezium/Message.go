package debeziumModels

import "github.com/google/uuid"

type DebeziumMessage struct {
	Payload DebeziumMessagePayload `json:"payload"`
}

type DebeziumMessagePayload struct {
	After DebeziumMessageNotionSave `json:"after"`
}

type DebeziumMessageNotionSave struct {
	Id       uuid.UUID `json:"id"`
	LastSave string    `json:"last_save"`
	Status   string    `json:"status"`
}

// Verifies if the the message is valid (it should at least have an Id and a LastSave)
func (d DebeziumMessage) IsValid() bool {
	return d.Payload.After.Id != uuid.Nil && d.Payload.After.LastSave != ""
}

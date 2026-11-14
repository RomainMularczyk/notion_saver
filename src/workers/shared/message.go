package sharedWorker

import (
	"encoding/json"
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/utils"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PageMessage struct {
	SaveId          uuid.UUID `json:"save_id"`
	PaginationIndex int       `json:"pagination_index"`
	NextCursor      *string   `json:"next_cursor"`
	RetryCount      int       `json:"retry_count"`
}

// Publishes a message in the queue with the pages under the next cursor
func PublishNextCursorMessage(
	ch *amqp.Channel,
	pageMessage PageMessage,
	server *utils.Server,
) {
	server.Logger.Info(
		"Publishing next cursor message",
		slog.String("NextCursor", *pageMessage.NextCursor),
	)
	messageBytes, err := json.Marshal(pageMessage)
	if err != nil {
		server.Logger.Error("Error when marshalling JSON")
		return
	}

	ch.Publish(
		"notion.pages",
		"notion.pages",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
			Headers:     amqp.Table{"x-retry-count": pageMessage.RetryCount},
		},
	)
}

// Publishes a message in the queue with the pages under the next cursor
func SendPageMessageToDlq(
	ch *amqp.Channel,
	pageMessage PageMessage,
	server *utils.Server,
) {
	server.Logger.Info("Sending page message to notion.pages.dlq.")
	messageBytes, err := json.Marshal(pageMessage)
	if err != nil {
		server.Logger.Error("Error when marshalling JSON.")
		return
	}

	ch.Publish(
		"notion.pages.dlq",
		"notion.pages.dlq",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Headers:     amqp.Table{"x-retry-count": pageMessage.RetryCount + 1},
			Body:        messageBytes,
		},
	)
}

// Retrieve the pages on the current cursor from the Notion API
func GetNextPages(
	ch *amqp.Channel,
	pageMessage *PageMessage,
	server *utils.Server,
) ([]saverModels.Page, *utils.NotionError) {
	notion := utils.NewNotion(server)
	nextNotionPages, notionErr := notion.SearchPages(
		utils.NotionFilterType("page"),
		pageMessage.NextCursor,
	)
	if notionErr != nil && notionErr.Kind == utils.NotionRateLimitError {
		server.Logger.Error(
			"Error when retrieving pages from Notion API.",
			slog.String("kind", string(notionErr.Kind)),
		)
		SendPageMessageToDlq(ch, *pageMessage, server)
		return nil, notionErr
	}

	if nextNotionPages.HasMore {
		nextPageMessage := PageMessage{
			NextCursor:      &nextNotionPages.NextCursor,
			SaveId:          pageMessage.SaveId,
			PaginationIndex: pageMessage.PaginationIndex + 1,
			RetryCount:      0,
		}
		PublishNextCursorMessage(ch, nextPageMessage, server)
	}

	newPages := utils.BuildSaverPages(nextNotionPages)
	return newPages, nil
}

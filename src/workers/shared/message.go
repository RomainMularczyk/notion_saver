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
		},
	)
}

// Publishes a message in the queue with the pages under the next cursor
func SendPageMessageToDql(
	ch *amqp.Channel,
	pageMessage PageMessage,
	server *utils.Server,
) {
	server.Logger.Info("Sending page message to notion.pages.dql.")
	messageBytes, err := json.Marshal(pageMessage)
	if err != nil {
		server.Logger.Error("Error when marshalling JSON.")
		return
	}

	ch.Publish(
		"notion.pages.dql",
		"notion.pages.dql",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
}

// Retrieve the pages on the current cursor from the Notion API
func GetNextPages(
	ch *amqp.Channel,
	pageMessage *PageMessage,
	server *utils.Server,
) ([]saverModels.Page, error) {
	notion := utils.NewNotion(server)
	nextNotionPages, err := notion.SearchPages(
		utils.NotionFilterType("page"),
		pageMessage.NextCursor,
	)
	if err != nil {
		notionError := notion.IsNotionError(err)
		if notionError.IsRateLimited {
			server.Logger.Error(
				"Error when retrieving pages from Notion API.",
				slog.String("error", err.Error()),
				slog.Bool("rate_limited", notionError.IsRateLimited),
				slog.Bool("client_error", notionError.IsClientError),
			)
			SendPageMessageToDql(ch, *pageMessage, server)
			return nil, err
		}
		return nil, err
	}

	if nextNotionPages.HasMore {
		nextPageMessage := PageMessage{
			NextCursor:      &nextNotionPages.NextCursor,
			SaveId:          pageMessage.SaveId,
			PaginationIndex: pageMessage.PaginationIndex + 1,
		}

		PublishNextCursorMessage(ch, nextPageMessage, server)
	}

	newPages := utils.BuildSaverPages(nextNotionPages)
	return newPages, nil
}

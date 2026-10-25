package blocksWorker

import (
	"encoding/json"
	"log/slog"
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/utils"
	blocksWorkerModels "notion_saver/src/workers/blocks/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Retrieve the blocks on the current cursor from the Notion API
// and publish the next cursor message (if there are more blocks
// to retrieve)
func (b *BlockWorker) GetNotionBlocks(
	ch *amqp.Channel,
	blockMessage blocksWorkerModels.BlockMessage,
) ([]saverModels.Block, error) {
	notion := utils.NewNotion(b.server)
	notionBlocks, err := notion.GetPageBlocks(
		blockMessage.PageId,
		blockMessage.NextCursor,
	)
	notionError := notion.IsNotionError(err)
	if err != nil && notionError.IsRateLimited {
		b.logger.Error(
			"Error when retrieving blocks from Notion API.",
			slog.String("error", err.Error()),
			slog.Bool("rate_limited", notionError.IsRateLimited),
			slog.Bool("client_error", notionError.IsClientError),
		)
		return nil, err
	}

	if err != nil {
		b.logger.Error(
			"Error when retrieving blocks from Notion API.",
			slog.String("error", err.Error()),
			slog.String("PageId", blockMessage.PageId.String()),
		)
		return nil, err
	}

	if notionBlocks.HasMore {
		nextBlockMessage := blocksWorkerModels.BlockMessage{
			NextCursor:      &notionBlocks.NextCursor,
			PageId:          blockMessage.PageId,
			PaginationIndex: blockMessage.PaginationIndex + 1,
			RetryCount:      0,
		}
		b.publishNextCursorMessage(ch, nextBlockMessage)
	}

	newBlocks := notionBlocks.ToSaverFormat(blockMessage.PageId)
	return newBlocks, nil
}

// Publish the next cursor message in the queue in order
// to retrieve the next cursor blocks
func (b *BlockWorker) publishNextCursorMessage(
	ch *amqp.Channel,
	blockMessage blocksWorkerModels.BlockMessage,
) {
	b.logger.Info(
		"Publishing next cursor message.",
		slog.String("page_id", blockMessage.PageId.String()),
	)
	messageBytes, err := json.Marshal(blockMessage)
	if err != nil {
		b.logger.Error("Error when marshalling JSON.")
		return
	}

	ch.Publish(
		"notion.blocks",
		"notion.blocks",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Headers:     amqp.Table{"x-retry-count": blockMessage.RetryCount},
			Body:        messageBytes,
		},
	)
}

// Publish the block message in the DLQ
func (b *BlockWorker) SendBlockMessageToDql(
	ch *amqp.Channel,
	blockMessage blocksWorkerModels.BlockMessage,
) {
	b.logger.Info("Sending block message to DQL.")
	messageBytes, err := json.Marshal(blockMessage)
	if err != nil {
		b.logger.Error("Error when marshalling JSON.")
		return
	}

	ch.Publish(
		"notion.blocks.dlq",
		"notion.blocks.dlq",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Headers:     amqp.Table{"x-retry-count": blockMessage.RetryCount + 1},
			Body:        messageBytes,
		},
	)
}

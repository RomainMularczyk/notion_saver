package blocksWorker

import (
	"encoding/json"
	"errors"
	"fmt"
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
	var notionErr *utils.NotionError
	if errors.As(err, &notionErr) {
		// TODO: implement exponential backoff
	}

	if err != nil {
		b.logger.Error(
			fmt.Sprintf("Error when retrieving blocks from Notion API: %v", err),
		)
		return nil, err
	}

	if notionBlocks.HasMore {
		nextBlockMessage := blocksWorkerModels.BlockMessage{
			NextCursor:      &notionBlocks.NextCursor,
			PageId:          blockMessage.PageId,
			PaginationIndex: blockMessage.PaginationIndex + 1,
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
		fmt.Sprintf("Publishing next cursor message: %v", utils.ToJson(blockMessage)),
	)
	messageBytes, err := json.Marshal(blockMessage)
	if err != nil {
		b.logger.Error(fmt.Sprintf("Error when marshalling JSON: %v", err))
		return
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

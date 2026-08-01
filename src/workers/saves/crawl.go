package savesWorker

import (
	saverModels "notion_saver/src/models/saver"
	"notion_saver/src/utils"

	amqp "github.com/rabbitmq/amqp091-go"
	sharedWorker "notion_saver/src/workers/shared"
)

// Retrieve the pages on the next cursor from the Notion API
// and publish the next cursor message (if there are more pages
// to retrieve)
func (w *SaveWorker) GetNextCursorPages(
	ch *amqp.Channel,
	pageMessage *sharedWorker.PageMessage,
) ([]saverModels.Page, error) {
	notion := utils.NewNotion(w.server)
	nextNotionPages, err := notion.SearchPages(
		utils.NotionFilterType("page"),
		pageMessage.NextCursor,
	)
	if err != nil {
		return nil, err
	}

	if nextNotionPages.HasMore {
		nextPageMessage := sharedWorker.PageMessage{
			NextCursor:      &nextNotionPages.NextCursor,
			SaveId:          pageMessage.SaveId,
			PaginationIndex: pageMessage.PaginationIndex + 1,
		}

		sharedWorker.PublishNextCursorMessage(ch, nextPageMessage, w.server)
	}

	newPages := utils.BuildSaverPages(nextNotionPages)
	return newPages, nil
}

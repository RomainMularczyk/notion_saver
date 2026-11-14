package sharedWorker

import (
	"notion_saver/src/queue"
	"notion_saver/src/utils"

	"github.com/rabbitmq/amqp091-go"
)

func NotionRetry(
	policy queue.RetryPolicy,
	d amqp091.Delivery,
	notionErr *utils.NotionError,
) {
	switch {
	case notionErr == nil:
		d.Ack(false)
	case !notionErr.IsRetryable():
		d.Ack()
	}
}

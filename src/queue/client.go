package queue

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"log/slog"
)

type Queue struct {
	Client *amqp.Connection
}

// Open the connection to RabbitMQ
func (q *Queue) Connect(
	host string,
	user string,
	password string,
	port string,
) {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port),
	)

	if err != nil {
		slog.Error(
			fmt.Sprintf(
				"Could not connect to RabbitMQ: %v",
				err,
			),
		)
		panic(err)
	}

	q.Client = conn
}

// Setup the RabbitMQ topology
func (q *Queue) SetupTopology() {
	ch, err := q.Client.Channel()
	defer ch.Close()
	if err != nil {
		slog.Error(
			fmt.Sprintf("Could not create channel: %v", err),
		)
		panic(err)
	}

	// EXCHANGES
	err = ch.ExchangeDeclare(
		"notion.pages",
		"direct",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	err = ch.ExchangeDeclare(
		"notion.pages.dlq",
		"direct",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	err = ch.ExchangeDeclare(
		"notion.blocks",
		"direct",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	err = ch.ExchangeDeclare(
		"notion.blocks.dlq",
		"direct",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	// QUEUES
	notionPagesQueue, err := ch.QueueDeclare(
		"notion.pages",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			amqp.QueueTypeArg:           amqp.QueueTypeQuorum,
			"x-dead-letter-exchange":    "notion.pages.dlq",
			"x-dead-letter-routing-key": "notion.pages.dlq",
		},
	)
	if err != nil {
		log.Panicln(err)
	}

	notionPagesDlq, err := ch.QueueDeclare(
		"notion.pages.dlq",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{amqp.QueueTypeArg: amqp.QueueTypeQuorum},
	)
	if err != nil {
		log.Panicln(err)
	}

	notionBlocksQueue, err := ch.QueueDeclare(
		"notion.blocks",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			amqp.QueueTypeArg:           amqp.QueueTypeQuorum,
			"x-dead-letter-exchange":    "notion.blocks.dlq",
			"x-dead-letter-routing-key": "notion.blocks.dlq",
		},
	)
	if err != nil {
		log.Panicln(err)
	}

	notionBlocksDlq, err := ch.QueueDeclare(
		"notion.blocks.dlq",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{amqp.QueueTypeArg: amqp.QueueTypeQuorum},
	)

	// BINDINGS
	err = ch.QueueBind(
		notionPagesQueue.Name,
		"notion.pages",
		"notion.pages",
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	err = ch.QueueBind(
		notionPagesDlq.Name,
		"notion.pages.dlq",
		"notion.pages.dlq",
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	err = ch.QueueBind(
		notionBlocksQueue.Name,
		"notion.blocks",
		"notion.blocks",
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	err = ch.QueueBind(
		notionBlocksDlq.Name,
		"notion.blocks.dlq",
		"notion.blocks.dlq",
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}
}

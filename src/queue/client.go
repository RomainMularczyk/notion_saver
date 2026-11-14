package queue

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
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
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port))

	if err != nil {
		log.Panicln("Could not connect to RabbitMQ.")
	}

	q.Client = conn
}

// Setup the RabbitMQ topology
func (q *Queue) SetupTopology() {
	ch, err := q.Client.Channel()
	if err != nil {
		log.Panicln("Could not create channel.")
	}
	defer ch.Close()

	// ----------------------------
	// --------- EXCHANGES --------
	// ----------------------------

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
		"notion.pages.dlx",
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
		"notion.blocks.dlx",
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

	// -------------------------
	// --------- QUEUES --------
	// -------------------------

	notionPagesQueue, err := ch.QueueDeclare(
		"notion.pages",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			amqp.QueueTypeArg:           amqp.QueueTypeQuorum,
			"x-dead-letter-exchange":    "notion.pages.dlx",
			"x-dead-letter-routing-key": "notion.pages.dlq",
			"x-delivery-limit":          int32(6),
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
			"x-dead-letter-exchange":    "notion.blocks.dlx",
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
	if err != nil {
		log.Panicln(err)
	}

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
		"notion.pages.dlx",
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
		"notion.blocks.dlx",
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(err)
	}

	// ---------------------------------
	// --------- RETRY TOPOLOGY --------
	// ---------------------------------

	savesRetry := NewRetryPolicy("notion.saves", DefaultRetryDelays)
	savesRetry.Declare(ch)

	pagesRetry := NewRetryPolicy("notion.pages", DefaultRetryDelays)
	pagesRetry.Declare(ch)

	blocksRetry := NewRetryPolicy("notion.blocks", DefaultRetryDelays)
	blocksRetry.Declare(ch)
}

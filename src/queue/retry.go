package queue

import (
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RetryPolicy struct {
	prefix      string // prefix will be "notion.saves" / "notion.pages" / "notion.blocks"
	delays      []time.Duration
	maxAttempts int
}

var DefaultRetryDelays = []time.Duration{
	5 * time.Second,
	30 * time.Second,
	2 * time.Minute,
	10 * time.Minute,
}

// Create a new retry policy
func NewRetryPolicy(prefix string, delays []time.Duration) RetryPolicy {
	return RetryPolicy{
		prefix:      prefix,
		delays:      delays,
		maxAttempts: len(delays),
	}
}

/*
	Create the exchange name.

For example, "notion.blocks.retry".
*/
func (p RetryPolicy) ExchangeName() string {
	return p.prefix + ".retry"
}

/*
	Create the queue name.

For example, "notion.blocks.retry.2".
*/
func (p RetryPolicy) queueName(attempt int) string {
	return fmt.Sprintf("%s.retry.%d", p.prefix, attempt)
}

/*
	Create the routing key.

For example, "retry.2".
*/
func (p RetryPolicy) routingKey(attempt int) string {
	return fmt.Sprintf("retry.%d", attempt)
}

func RetryCount(d amqp091.Delivery) int {
	var retryCount int
	switch count := d.Headers["x-retry-count"].(type) {
	case int32:
		retryCount = int(count)
	case int64:
		retryCount = int(count)
	case float32:
		retryCount = int(count)
	case float64:
		retryCount = int(count)
	default:
		retryCount = 0
	}
	if retryCount < 0 {
		return 0
	} else {
		return retryCount
	}
}

func (p RetryPolicy) NextAttempt(d amqp091.Delivery) int {
	next := RetryCount(d) + 1
	return next
}

/*
	Schedule the message in the corresponding retry queue. This is decided based on the

number of attempts value.
*/
func (p RetryPolicy) Schedule(
	ch *amqp091.Channel,
	d amqp091.Delivery,
) error {
	attempt := p.NextAttempt(d)
	if attempt > p.maxAttempts {
		return fmt.Errorf("Max attempt reached.")
	}

	// NOTE: we need to copy headers to preserve OTel trace correlation headers
	headers := make(amqp091.Table, len(d.Headers)+1)
	for k, v := range d.Headers {
		headers[k] = v
	}
	headers["x-retry-count"] = int32(attempt)

	err := ch.Publish(
		p.ExchangeName(),
		p.routingKey(attempt),
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Headers:      headers,
			Body:         d.Body,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"Publish retry attempt %d. Error: %w",
			attempt,
			err,
		)
	}
	return nil
}

// Declare both the exchange, the queue and the queue binding to RabbitMQ
func (p RetryPolicy) Declare(ch *amqp091.Channel) {
	err := ch.ExchangeDeclare(
		p.ExchangeName(),
		"direct",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Panicln(fmt.Sprintf(
			"Declaring exchange %s failed with error: %w",
			p.ExchangeName(),
			err,
		))
	}

	for i, delay := range p.delays {
		attempt := i + 1
		queueName := p.queueName(attempt)

		_, err := ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			amqp091.Table{
				"x-message-ttl":             int32(delay.Milliseconds()),
				"x-dead-letter-exchange":    p.ExchangeName(),
				"x-dead-letter-routing-key": p.prefix,
			},
		)
		if err != nil {
			log.Panicln(fmt.Sprintf(
				"Could not create queue %s",
				queueName,
			))
		}

		err = ch.QueueBind(
			queueName,
			p.routingKey(attempt),
			p.ExchangeName(),
			false,
			nil,
		)
		if err != nil {
			log.Panicln(fmt.Sprintf(
				"Could not bind queue %s with exchange %s.",
				queueName,
				p.ExchangeName(),
			))
		}
	}
}

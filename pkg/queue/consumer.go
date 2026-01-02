package queue

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	queueClient *QueueClient
}

func NewConsumer(queueName string) *Consumer {
	queueClient := NewQueueClient(queueName)
	return &Consumer{
		queueClient: queueClient,
	}
}

func (c *Consumer) Close() {
	c.queueClient.Close()
}

func (c *Consumer) ConsumeMessages(ctx context.Context) (<-chan amqp091.Delivery, error) {
	msgs, err := c.queueClient.queue.ConsumeWithContext(
		ctx,
		c.queueClient.queueName,
		"",
		false, // no auto-ack
		false, // not exclusive
		false, // no local
		false, // no wait
		nil,
	)
	if err != nil {
		slog.Error("[queue.Consumer.ConsumeMessages] Failed to consume messages", "error", err)
		return nil, err
	}

	return msgs, nil
}

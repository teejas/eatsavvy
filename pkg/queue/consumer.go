package queue

import (
	"context"
	"fmt"
	"time"
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

func (c *Consumer) ConsumeMessages() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msgs, err := c.queueClient.queue.ConsumeWithContext(
		ctx,
		"",
		c.queueClient.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		fmt.Println(string(msg.Body))
		msg.Ack(false)
	}
	return nil
}

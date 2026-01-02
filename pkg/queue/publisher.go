package queue

import (
	"context"
	"log/slog"
	"time"

	"eatsavvy/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	queueClient *QueueClient
}

func NewPublisher(queueName string) *Publisher {
	queueClient := NewQueueClient(queueName)
	return &Publisher{
		queueClient: queueClient,
	}
}

func (p *Publisher) Close() {
	p.queueClient.Close()
}

func (p *Publisher) PublishMessage(body interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bodyBytes, err := utils.ToBytes(body)
	slog.Info("[queue.Publisher.PublishMessage] Publishing message", "body", body)
	if err != nil {
		slog.Error("[queue.Publisher.PublishMessage] Failed to convert body to bytes", "error", err)
		return err
	}
	err = p.queueClient.queue.PublishWithContext(
		ctx,
		"",
		p.queueClient.queueName,
		false,
		false,
		amqp.Publishing{
			Body: bodyBytes,
		},
	)
	if err != nil {
		slog.Error("[queue.Publisher.PublishMessage] Failed to publish message", "error", err)
		return err
	}
	return nil
}

func (p *Publisher) PublishDelayedMessage(body interface{}, delay time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bodyBytes, err := utils.ToBytes(body)
	if err != nil {
		slog.Error("[queue.Publisher.PublishDelayedMessage] Failed to convert body to bytes", "error", err)
		return err
	}
	err = p.queueClient.queue.PublishWithContext(
		ctx,
		"delayed-exchange",
		p.queueClient.queueName,
		false,
		false,
		amqp.Publishing{
			Body: bodyBytes,
			Headers: amqp.Table{
				"x-delay": delay.Milliseconds(),
			},
		},
	)
	slog.Info("[queue.Publisher.PublishDelayedMessage] Published message with delay", "delay", delay.Milliseconds())
	if err != nil {
		slog.Error("[queue.Publisher.PublishDelayedMessage] Failed to publish message", "error", err)
		return err
	}
	return nil
}

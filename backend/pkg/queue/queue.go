package queue

import (
	"log/slog"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueClient struct {
	queue     *amqp.Channel
	conn      *amqp.Connection
	queueName string
}

func NewQueueClient(queueName string) *QueueClient {
	conn, err := amqp.Dial("amqp://" +
		os.Getenv("RABBITMQ_USER") + ":" +
		os.Getenv("RABBITMQ_PASSWORD") + "@" +
		os.Getenv("RABBITMQ_HOST") + ":" +
		os.Getenv("RABBITMQ_PORT") + "/")

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("[queue.NewQueueClient] Failed to open a channel", "error", err)
	}

	// Declare a delayed message exchange
	err = ch.ExchangeDeclare(
		"delayed-exchange", "x-delayed-message", true, false, false, false, amqp.Table{
			"x-delayed-type": "direct",
		},
	)

	if err != nil {
		slog.Error("[queue.NewQueueClient] Failed to declare delayed exchange", "error", err)
	}

	_, err = ch.QueueDeclare(
		queueName, false, false, false, false, nil,
	)
	if err != nil {
		slog.Error("[queue.NewQueueClient] Failed to create queue", "error", err)
	}

	// Bind the queue to the delayed exchange
	err = ch.QueueBind(
		queueName,          // queue name
		queueName,          // routing key (same as queue name)
		"delayed-exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		slog.Error("[queue.NewQueueClient] Failed to bind queue to delayed exchange", "error", err)
	}

	return &QueueClient{
		conn:      conn,
		queue:     ch,
		queueName: queueName,
	}
}

func (qc *QueueClient) Close() {
	slog.Info("[queue.Close] Closing connection and channel", "queueName", qc.queueName)
	err := qc.conn.Close()
	if err != nil {
		slog.Error("[queue.Close] Failed to close connection", "error", err)
	}
	err = qc.queue.Close()
	if err != nil {
		slog.Error("[queue.Close] Failed to close channel", "error", err)
	}
	slog.Info("[queue.Close] Closed connection and channel", "queueName", qc.queueName)
}

package worker

import (
	"eatsavvy/pkg/queue"
	"log/slog"
)

type Worker struct {
	consumer *queue.Consumer
}

func NewWorker() *Worker {
	consumer := queue.NewConsumer("enrich_restaurant_details")
	return &Worker{
		consumer: consumer,
	}
}

func (w *Worker) Start() {
	slog.Info("[worker.Start] Starting worker")

	err := w.consumer.ConsumeMessages()
	if err != nil {
		slog.Error("[worker.Start] Failed to consume messages", "error", err)
		return
	}
}

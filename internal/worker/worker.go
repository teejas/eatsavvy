package worker

import (
	"context"
	"eatsavvy/pkg/db"
	"eatsavvy/pkg/places"
	"eatsavvy/pkg/queue"
	"eatsavvy/pkg/utils"
	"eatsavvy/pkg/vapi"

	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	consumer   *queue.Consumer
	publisher  *queue.Publisher
	vapiClient *vapi.VapiClient
	dbClient   *db.DatabaseClient
}

func NewWorker() *Worker {
	consumer := queue.NewConsumer("enrich_restaurant_details")
	publisher := queue.NewPublisher("enrich_restaurant_details")
	vapiClient := vapi.NewVapiClient()
	dbClient := db.NewDatabaseClient()
	return &Worker{
		consumer:   consumer,
		publisher:  publisher,
		vapiClient: vapiClient,
		dbClient:   dbClient,
	}
}

func (w *Worker) Close() {
	w.consumer.Close()
	w.publisher.Close()
	w.dbClient.Close()
}

func (w *Worker) Start() {
	slog.Info("[worker.Start] Starting worker")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer w.Close()

	msgs, err := w.consumer.ConsumeMessages(ctx)
	if err != nil {
		slog.Error("[worker.Start] Failed to consume messages", "error", err)
		return
	}

	forever := make(chan struct{})

	go func() {
		for msg := range msgs {
			restaurant, err := w.processMessage(msg)
			if err != nil {
				slog.Error("[worker.processMessages] Failed to process message", "error", err)
				if restaurant.Id != "" {
					err = w.handleFailure(restaurant.Id)
					if err != nil {
						slog.Error("[worker.processMessages] Failed to handle failure", "error", err)
					}
				} else {
					slog.Error("[worker.processMessages] Failed to get restaurant ID", "error", err)
				}
				continue
			}
			msg.Ack(false)
		}
		slog.Info("[worker.processMessages] Waiting for more messages...")
	}()

	slog.Info("[worker.Start] Ready to receive messages")

	<-forever
}

func (w *Worker) processMessage(msg amqp091.Delivery) (places.Restaurant, error) {
	var restaurant places.Restaurant
	err := utils.FromBytes(msg.Body, &restaurant)
	if err != nil {
		slog.Error("[worker.processMessage] Failed to unmarshal message", "error", err)
		return places.Restaurant{}, err
	}
	slog.Info("[worker.processMessage] Processing message", "message", restaurant)
	now := time.Now().UTC()
	currentDay := int(now.Weekday())
	currentHour := now.Hour()
	openNow := false
	for _, openHour := range restaurant.OpenHours {
		if (openHour.Open.Weekday == currentDay && openHour.Open.Hour < currentHour) &&
			(openHour.Close.Hour > currentHour || openHour.Close.Weekday > currentDay) {
			openNow = true
			slog.Info("[worker.processMessage] Restaurant is open", "restaurant", restaurant.Name)
			vapiResponse, err := w.vapiClient.CreateCall(restaurant)
			if err != nil {
				slog.Error("[worker.processMessage] Failed to make Vapi phone call", "error", err)
				return restaurant, err
			}
			slog.Info("[worker.processMessage] Vapi phone call made", "callId", vapiResponse.Id)
			_, err = w.dbClient.Db.Exec(w.dbClient.Ctx,
				`UPDATE public.restaurants SET enrichment_status = $1 WHERE places_id = $2`,
				places.EnrichmentStatusInProgress, restaurant.Id,
			)
			if err != nil {
				slog.Error("[worker.processMessage] Failed to update enrichment status", "error", err)
				return restaurant, err
			}
			_, err = w.dbClient.Db.Exec(w.dbClient.Ctx,
				`INSERT INTO public.calls (places_id, vapi_call_id, call_status) VALUES ($1, $2, $3)`,
				restaurant.Id, vapiResponse.Id, "initiated",
			)
			if err != nil {
				slog.Error("[worker.processMessage] Failed to update Vapi call ID", "error", err)
				return restaurant, err
			}
			return restaurant, nil
		}
	}
	if !openNow {
		slog.Info("[worker.processMessage] Restaurant is closed", "restaurant", restaurant.Name)
		err := w.publisher.PublishDelayedMessage(restaurant, 1*time.Hour)
		if err != nil {
			slog.Error("[worker.processMessage] Failed to publish message", "error", err)
			return restaurant, err
		}
		slog.Info("[worker.processMessage] Published message with delay", "delay", 1*time.Hour)
	}
	return restaurant, nil
}

func (w *Worker) handleFailure(restaurantId string) error {
	_, err := w.dbClient.Db.Exec(w.dbClient.Ctx,
		`UPDATE public.restaurants SET enrichment_status = $1 WHERE places_id = $2`,
		places.EnrichmentStatusFailed, restaurantId,
	)
	if err != nil {
		slog.Error("[worker.handleFailure] Failed to update enrichment status", "error", err)
		return err
	}
	slog.Info("[worker.processMessage] Updated enrichment status to failed", "places_id", restaurantId)
	return nil
}

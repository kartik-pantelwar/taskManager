package subscriber

import (
	"context"
	"encoding/json"
	"log"
	"notificationservice/src/internal/adaptors/repository"
	"notificationservice/src/internal/core/task"
	"notificationservice/src/internal/usecase"
)

type EventSubscriber struct {
	notificationRepo    *repository.NotificationRepository
	notificationUseCase *usecase.NotificationUseCase
}

func NewEventSubscriber(repo *repository.NotificationRepository, uc *usecase.NotificationUseCase) *EventSubscriber {
	return &EventSubscriber{
		notificationRepo:    repo,
		notificationUseCase: uc,
	}
}

func (s *EventSubscriber) StartListening(ctx context.Context, channel string) {
	log.Println("Starting notification service event listener...")

	// Subscribe to task events
	pubsub := s.notificationRepo.Subscribe(ctx, channel)
	defer pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("Event subscriber shutting down...")
			return
		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				continue
			}

			var event task.TaskEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				continue
			}

			log.Printf("Received event: %s for task %d", event.EventType, event.TaskID)

			// Process the event and store notification
			if err := s.notificationUseCase.ProcessTaskEvent(ctx, event); err != nil {
				log.Printf("Failed to process event: %v", err)
			}
		}
	}
}

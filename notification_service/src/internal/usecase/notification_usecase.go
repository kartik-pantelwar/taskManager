package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notificationservice/src/internal/adaptors/redis"
	"notificationservice/src/internal/core/notification"
	"notificationservice/src/internal/core/task"
	"time"

	"github.com/google/uuid"
)

type NotificationUseCase struct {
	redisClient *redis.RedisClient
}

func NewNotificationUseCase(redisClient *redis.RedisClient) *NotificationUseCase {
	return &NotificationUseCase{
		redisClient: redisClient,
	}
}

func (uc *NotificationUseCase) ProcessTaskEvent(ctx context.Context, event task.TaskEvent) error {
	// Create notification for assigned user
	notif := notification.Notification{
		ID:         uuid.New().String(),
		TaskID:     event.TaskID,
		Action:     event.EventType,
		TaskName:   event.TaskName,
		UserID:     event.AssignedTo, // User who will receive notification
		AssignedBy: event.AssignedBy, // User who assigned the task (logged-in user)
		Message:    uc.generateMessage(event.EventType, event.TaskName, event.AssignedTo),
		Timestamp:  time.Now(),
	}

	// Store notification in Redis
	notificationJSON, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	notificationKey := fmt.Sprintf("notification:%s", notif.ID)
	err = uc.redisClient.StoreNotification(ctx, notificationKey, notificationJSON)
	if err != nil {
		return fmt.Errorf("failed to store notification: %v", err)
	}

	log.Printf("Notification stored for user %d: %s (assigned by user %d)", notif.UserID, notif.Message, notif.AssignedBy)
	return nil
}

func (uc *NotificationUseCase) GetMostRecentNotification(ctx context.Context) (*notification.Notification, error) {
	// Get all notification keys
	keys, err := uc.redisClient.GetAllNotificationKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification keys: %v", err)
	}

	if len(keys) == 0 {
		return nil, nil
	}

	// Just return the first notification found (no complex sorting)
	for _, key := range keys {
		notificationJSON, err := uc.redisClient.GetNotification(ctx, key)
		if err != nil {
			continue
		}

		var notif notification.Notification
		err = json.Unmarshal(notificationJSON, &notif)
		if err != nil {
			continue
		}

		return &notif, nil
	}

	return nil, nil
}

func (uc *NotificationUseCase) GetAllNotifications(ctx context.Context, limit int) ([]notification.Notification, error) {
	// Get all notification keys
	keys, err := uc.redisClient.GetAllNotificationKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification keys: %v", err)
	}

	var notifications []notification.Notification

	// Get all notifications (no complex filtering or sorting)
	for _, key := range keys {
		notificationJSON, err := uc.redisClient.GetNotification(ctx, key)
		if err != nil {
			continue
		}

		var notif notification.Notification
		err = json.Unmarshal(notificationJSON, &notif)
		if err != nil {
			continue
		}

		notifications = append(notifications, notif)
	}

	return notifications, nil
}

func (uc *NotificationUseCase) GetMyNotifications(ctx context.Context, assignedBy int, limit int) ([]notification.Notification, error) {
	// Get all notification keys
	keys, err := uc.redisClient.GetAllNotificationKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification keys: %v", err)
	}

	var notifications []notification.Notification

	// Simple filter by assignedBy (no complex sorting)
	for _, key := range keys {
		notificationJSON, err := uc.redisClient.GetNotification(ctx, key)
		if err != nil {
			continue
		}

		var notif notification.Notification
		err = json.Unmarshal(notificationJSON, &notif)
		if err != nil {
			continue
		}

		// Simple filter - add if assigned by matches
		if notif.AssignedBy == assignedBy {
			notifications = append(notifications, notif)
		}
	}

	return notifications, nil
}

func (uc *NotificationUseCase) GetMyRecentNotification(ctx context.Context, assignedBy int) (*notification.Notification, error) {
	// Get all notification keys
	keys, err := uc.redisClient.GetAllNotificationKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification keys: %v", err)
	}

	// Simple search - return first match found (no complex sorting)
	for _, key := range keys {
		notificationJSON, err := uc.redisClient.GetNotification(ctx, key)
		if err != nil {
			continue
		}

		var notif notification.Notification
		err = json.Unmarshal(notificationJSON, &notif)
		if err != nil {
			continue
		}

		// Return first notification where assignedBy matches
		if notif.AssignedBy == assignedBy {
			return &notif, nil
		}
	}

	return nil, nil
}

func (uc *NotificationUseCase) generateMessage(action, taskName string, assignedTo int) string {
	switch action {
	case "task_created":
		return fmt.Sprintf("Task '%s' assigned to user %d", taskName, assignedTo)
	case "task_updated":
		return fmt.Sprintf("Task '%s' updated (assigned to user %d)", taskName, assignedTo)
	case "task_deleted":
		return fmt.Sprintf("Task '%s' deleted (was assigned to user %d)", taskName, assignedTo)
	default:
		return fmt.Sprintf("Action '%s' performed on task '%s' (assigned to user %d)", action, taskName, assignedTo)
	}
}

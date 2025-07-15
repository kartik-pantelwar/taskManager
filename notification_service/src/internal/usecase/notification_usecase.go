package usecase

import (
	"context"
	"fmt"
	"log"
	"notificationservice/src/internal/adaptors/repository"
	"notificationservice/src/internal/core/notification"
	"notificationservice/src/internal/core/task"
	"time"

	"github.com/google/uuid"
)

type NotificationUseCase struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationUseCase(repo *repository.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo: repo,
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
		Message:    uc.generateMessage(event.EventType, event.TaskName, event.AssignedBy, event.AssignedTo),
		Timestamp:  time.Now(),
	}

	// Store notification in Redis
	err := uc.notificationRepo.StoreNotification(ctx, notif)
	if err != nil {
		return fmt.Errorf("failed to store notification: %v", err)
	}

	log.Printf("Notification stored for user %d: %s (assigned by user %d)", notif.UserID, notif.Message, notif.AssignedBy)
	return nil
}

func (uc *NotificationUseCase) GetMostRecentNotification(ctx context.Context) (*notification.Notification, error) {
	// Get the most recent notification globally
	recentNotification, err := uc.notificationRepo.GetMostRecentNotification(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get most recent notification: %v", err)
	}

	return recentNotification, nil
}

func (uc *NotificationUseCase) GetAllNotifications(ctx context.Context, limit int) ([]notification.Notification, error) {
	if limit <= 0 || limit > 200 {
		limit = 50 // Default limit
	}

	return uc.notificationRepo.GetAllNotifications(ctx, limit)
}

func (uc *NotificationUseCase) generateMessage(action, taskName string, assignedBy, assignedTo int) string {
	switch action {
	case "task_created":
		return fmt.Sprintf("User %d assigned you a new task: '%s'", assignedBy, taskName)
	case "task_updated":
		return fmt.Sprintf("User %d updated task: '%s'", assignedBy, taskName)
	case "task_deleted":
		return fmt.Sprintf("User %d deleted task: '%s'", assignedBy, taskName)
	default:
		return fmt.Sprintf("User %d performed action '%s' on task: '%s'", assignedBy, action, taskName)
	}
}

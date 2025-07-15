package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"task_service/src/internal/core/task"
	"time"

	"github.com/redis/go-redis/v9"
)

type NotificationService struct {
	redisClient *redis.Client
}

// Don't used
type TaskNotification struct {
	ID        string    `json:"id"`
	TaskID    int       `json:"task_id"`
	Action    string    `json:"action"` // "created", "updated", "deleted"
	TaskName  string    `json:"task_name"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Notification struct {
	ID        string    `json:"id"`
	TaskID    int       `json:"task_id"`
	Action    string    `json:"action"`
	TaskName  string    `json:"task_name"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewNotificationService(redisClient *redis.Client) *NotificationService {
	return &NotificationService{
		redisClient: redisClient,
	}
}

func (n *NotificationService) PublishTaskNotification(ctx context.Context, taskData task.Task, action string, userID int) error {
	notification := TaskNotification{
		ID:        fmt.Sprintf("%s_%d_%d", action, taskData.Id, time.Now().Unix()),
		TaskID:    taskData.Id,
		Action:    action,
		TaskName:  taskData.Name,
		UserID:    userID,
		Message:   n.generateMessage(action, taskData.Name),
		Timestamp: time.Now(),
	}

	// Convert to JSON
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	// Publish to Redis channel
	channel := "task_notifications" // !to be changed
	err = n.redisClient.Publish(ctx, channel, notificationJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to publish notification: %v", err)
	}

	// Store notification in Redis (for retrieval later)
	notificationKey := fmt.Sprintf("notification:%s", notification.ID)
	err = n.redisClient.Set(ctx, notificationKey, notificationJSON, 24*time.Hour).Err() // Store for 24 hours
	if err != nil {
		return fmt.Errorf("failed to store notification: %v", err)
	}

	// Add to user's notification list
	userNotificationsKey := fmt.Sprintf("user_notifications:%d", userID)
	err = n.redisClient.LPush(ctx, userNotificationsKey, notification.ID).Err()
	if err != nil {
		return fmt.Errorf("failed to add to user notifications: %v", err)
	}

	// Keep only last 50 notifications per user
	n.redisClient.LTrim(ctx, userNotificationsKey, 0, 49)

	return nil
}

func (n *NotificationService) generateMessage(action, taskName string) string {
	switch action {
	case "created":
		return fmt.Sprintf("Task '%s' has been created", taskName)
	case "updated":
		return fmt.Sprintf("Task '%s' has been updated", taskName)
	case "deleted":
		return fmt.Sprintf("Task '%s' has been deleted", taskName)
	default:
		return fmt.Sprintf("Task '%s' action: %s", taskName, action)
	}
}

func (n *NotificationService) GetUserNotifications(ctx context.Context, userID int, limit int) ([]TaskNotification, error) {
	userNotificationsKey := fmt.Sprintf("user_notifications:%d", userID)

	// Get notification IDs
	notificationIDs, err := n.redisClient.LRange(ctx, userNotificationsKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get notification IDs: %v", err)
	}

	var notifications []TaskNotification
	for _, id := range notificationIDs {
		notificationKey := fmt.Sprintf("notification:%s", id)
		notificationJSON, err := n.redisClient.Get(ctx, notificationKey).Result()
		if err != nil {
			continue // Skip if notification not found (might be expired)
		}

		var notification TaskNotification
		err = json.Unmarshal([]byte(notificationJSON), &notification)
		if err != nil {
			continue // Skip invalid JSON
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// PublishEvent publishes an event to a Redis channel
func (n *NotificationService) PublishEvent(ctx context.Context, channel string, eventJSON []byte) error {
	return n.redisClient.Publish(ctx, channel, eventJSON).Err()
}

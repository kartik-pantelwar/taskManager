package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"notificationservice/src/internal/core/notification"
	"time"

	"github.com/redis/go-redis/v9"
)

type NotificationRepository struct {
	redisClient *redis.Client
}

func NewNotificationRepository(redisClient *redis.Client) *NotificationRepository {
	return &NotificationRepository{
		redisClient: redisClient,
	}
}

func (r *NotificationRepository) StoreNotification(ctx context.Context, notif notification.Notification) error {
	// Convert to JSON
	notificationJSON, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	// Store individual notification (24-hour expiry)
	notificationKey := fmt.Sprintf("notification:%s", notif.ID)
	err = r.redisClient.Set(ctx, notificationKey, notificationJSON, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store notification: %v", err)
	}

	// Add to user's notification list (newest first)
	userNotificationsKey := fmt.Sprintf("user_notifications:%d", notif.UserID)
	err = r.redisClient.LPush(ctx, userNotificationsKey, notif.ID).Err()
	if err != nil {
		return fmt.Errorf("failed to add to user notifications: %v", err)
	}

	// Keep only last 50 notifications per user
	err = r.redisClient.LTrim(ctx, userNotificationsKey, 0, 49).Err()
	if err != nil {
		return fmt.Errorf("failed to trim user notifications: %v", err)
	}

	return nil
}

func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID int, limit int) ([]notification.Notification, error) {
	userNotificationsKey := fmt.Sprintf("user_notifications:%d", userID)

	// Get notification IDs
	notificationIDs, err := r.redisClient.LRange(ctx, userNotificationsKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get notification IDs: %v", err)
	}

	var notifications []notification.Notification
	for _, id := range notificationIDs {
		notificationKey := fmt.Sprintf("notification:%s", id)
		notificationJSON, err := r.redisClient.Get(ctx, notificationKey).Result()
		if err != nil {
			continue // Skip if notification not found (might be expired)
		}

		var notif notification.Notification
		err = json.Unmarshal([]byte(notificationJSON), &notif)
		if err != nil {
			continue // Skip invalid JSON
		}

		notifications = append(notifications, notif)
	}

	return notifications, nil
}

func (r *NotificationRepository) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.redisClient.Subscribe(ctx, channel)
}

func (r *NotificationRepository) GetMostRecentNotification(ctx context.Context) (*notification.Notification, error) {
	// Get the most recent notification by finding the latest added notification across all users
	// We can scan through all notification keys and find the one with the latest timestamp

	// Use SCAN to get all notification keys
	var cursor uint64
	var allNotificationKeys []string

	for {
		keys, newCursor, err := r.redisClient.Scan(ctx, cursor, "notification:*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification keys: %v", err)
		}

		allNotificationKeys = append(allNotificationKeys, keys...)
		cursor = newCursor

		if cursor == 0 {
			break
		}
	}

	if len(allNotificationKeys) == 0 {
		return nil, nil // No notifications found
	}

	var mostRecentNotification *notification.Notification
	var latestTimestamp time.Time

	// Check each notification and find the most recent one
	for _, key := range allNotificationKeys {
		notificationJSON, err := r.redisClient.Get(ctx, key).Result()
		if err != nil {
			continue // Skip if notification not found (might be expired)
		}

		var notif notification.Notification
		err = json.Unmarshal([]byte(notificationJSON), &notif)
		if err != nil {
			continue // Skip invalid JSON
		}

		if mostRecentNotification == nil || notif.Timestamp.After(latestTimestamp) {
			mostRecentNotification = &notif
			latestTimestamp = notif.Timestamp
		}
	}

	return mostRecentNotification, nil
}

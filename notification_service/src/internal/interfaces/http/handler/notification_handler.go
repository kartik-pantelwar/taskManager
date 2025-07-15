package handler

import (
	"encoding/json"
	"net/http"
	"notificationservice/src/internal/core/notification"
	"notificationservice/src/internal/usecase"
	"strconv"
)

type NotificationHandler struct {
	notificationUseCase *usecase.NotificationUseCase
}

func NewNotificationHandler(uc *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: uc,
	}
}

func (h *NotificationHandler) GetRecentNotification(w http.ResponseWriter, r *http.Request) {
	// Check if x-userId header is present for user-specific recent notification
	userIDHeader := r.Header.Get("x-userId")

	var notification *notification.Notification
	var err error

	if userIDHeader != "" {
		// Get recent notification for specific user (where they are assigned_by)
		userID, parseErr := strconv.Atoi(userIDHeader)
		if parseErr == nil {
			notification, err = h.notificationUseCase.GetMyRecentNotification(r.Context(), userID)
		} else {
			// Fall back to global recent if invalid user ID
			notification, err = h.notificationUseCase.GetMostRecentNotification(r.Context())
		}
	} else {
		// Get global recent notification
		notification, err = h.notificationUseCase.GetMostRecentNotification(r.Context())
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to retrieve recent notification",
		})
		return
	}

	// Return single recent notification or null if none exists
	if notification == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"recent_notification": nil,
			"message":             "no recent notifications found",
		})
		return
	}

	// Transform notification to show assigned_to instead of assigned_by
	transformedNotification := map[string]interface{}{
		"id":          notification.ID,
		"task_id":     notification.TaskID,
		"action":      notification.Action,
		"task_name":   notification.TaskName,
		"assigned_to": notification.UserID, // Show who the task was assigned TO
		"message":     notification.Message,
		"timestamp":   notification.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recent_notification": transformedNotification,
	})
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	// Check if x-userId header is present for user-specific notifications
	userIDHeader := r.Header.Get("x-userId")

	// Get limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default for all notifications
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	var notifications []notification.Notification
	var err error

	if userIDHeader != "" {
		// Get notifications where logged-in user is assigned_by
		userID, parseErr := strconv.Atoi(userIDHeader)
		if parseErr == nil {
			notifications, err = h.notificationUseCase.GetMyNotifications(r.Context(), userID, limit)
		} else {
			// Fall back to all notifications if invalid user ID
			notifications, err = h.notificationUseCase.GetAllNotifications(r.Context(), limit)
		}
	} else {
		// Get all notifications
		notifications, err = h.notificationUseCase.GetAllNotifications(r.Context(), limit)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to retrieve notifications",
		})
		return
	}

	// Transform notifications to show assigned_to instead of assigned_by
	transformedNotifications := make([]map[string]interface{}, len(notifications))
	for i, notif := range notifications {
		transformedNotifications[i] = map[string]interface{}{
			"id":          notif.ID,
			"task_id":     notif.TaskID,
			"action":      notif.Action,
			"task_name":   notif.TaskName,
			"assigned_to": notif.UserID, // Show who the task was assigned TO
			"message":     notif.Message,
			"timestamp":   notif.Timestamp,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notifications": transformedNotifications,
		"count":         len(transformedNotifications),
	})
}

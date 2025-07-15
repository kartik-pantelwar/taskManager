package handler

import (
	"encoding/json"
	"net/http"
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
	// Get the most recent notification globally (from all users)
	notification, err := h.notificationUseCase.GetMostRecentNotification(r.Context())
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recent_notification": notification,
	})
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	// Get limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default for all notifications
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	notifications, err := h.notificationUseCase.GetAllNotifications(r.Context(), limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to retrieve notifications",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notifications": notifications,
		"count":         len(notifications),
	})
}

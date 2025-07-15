package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NotificationHandler struct {
	notificationServiceURL string
}

func NewNotificationHandler(notificationServiceURL string) *NotificationHandler {
	return &NotificationHandler{
		notificationServiceURL: notificationServiceURL,
	}
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from context (now as int)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "user not found in context",
		})
		return
	}

	// Get limit from query parameter
	limit := "20" // default
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limit = limitParam
	}

	// Call notification service
	notificationURL := fmt.Sprintf("%s/api/v1/notifications/user/%d?limit=%s",
		h.notificationServiceURL, userID, limit)

	resp, err := http.Get(notificationURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to retrieve notifications from notification service",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "notification service returned error",
		})
		return
	}

	// Forward the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to read notification service response",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

package handler

import (
	"net/http"
	"notificationservice/src/internal/core/notification"
	"notificationservice/src/internal/usecase"
	errorhandling "notificationservice/src/pkg/error_handling"
	pkgresponse "notificationservice/src/pkg/response"
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
		errorhandling.HandleError(w, "Failed to Retrieve Recent Notification", http.StatusInternalServerError)
		return
	}

	// Return single recent notification or null if none exists
	if notification == nil {
		response := pkgresponse.StandardResponse{
			Status:  "SUCCESS",
			Message: "No Recent Notifications Found",
			Data:    nil,
		}
		pkgresponse.WriteResponse(w, http.StatusOK, response)
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

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Recent Notification Retrieved Successfully",
		Data:    transformedNotification,
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
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
		errorhandling.HandleError(w, "Failed to Retrieve Notifications", http.StatusInternalServerError)
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

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Notifications Retrieved Successfully",
		Data: map[string]interface{}{
			"notifications": transformedNotifications,
			"count":         len(transformedNotifications),
			// "user_filter":   userIDHeader != "",
		},
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

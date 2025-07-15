package routes

import (
	"net/http"
	"notificationservice/src/internal/interfaces/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRoutes(notificationHandler *handler.NotificationHandler) http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)

	// Health check
	router.Get("/health", notificationHandler.HealthCheck)

	// Notification routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/notifications", func(r chi.Router) {
			r.Get("/recent", notificationHandler.GetRecentNotification)
			r.Get("/user/{userID}", notificationHandler.GetUserNotifications)
		})
	})

	return router
}

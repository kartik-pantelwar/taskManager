package routes

import (
	"context"
	"net/http"
	"strconv"

	taskhandler "task_service/src/internal/interfaces/input/api/rest/handler"
	pb "task_service/src/internal/interfaces/input/grpc/generated/generated"

	"github.com/go-chi/chi/v5"
)

func SessionAuthMiddleware(grpcClient pb.SessionValidatorClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract "sess" cookie (refresh token/session ID)
			cookie, err := r.Cookie("sess")
			if err != nil {
				http.Error(w, "session cookie is missing", http.StatusUnauthorized)
				return
			}
			sessionID := cookie.Value

			// gRPC call to user service for session validation
			resp, err := grpcClient.ValidateSession(context.Background(), &pb.ValidateSessionRequest{
				SessionId: sessionID,
			})
			if err != nil || !resp.Valid {
				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}

			// Converting user id into int
			userID, err := strconv.Atoi(resp.UserId)
			if err != nil {
				http.Error(w, "invalid user ID from session", http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func InitRoutes(taskHandler *taskhandler.TaskHandler, notificationHandler *taskhandler.NotificationHandler, grpcClient pb.SessionValidatorClient) http.Handler {
	router := chi.NewRouter()

	router.Route("/v1/tasks", func(r chi.Router) {
		r.Use(SessionAuthMiddleware(grpcClient))
		r.Post("/create", taskHandler.Create)
		r.Put("/update", taskHandler.Update)
		r.Delete("/{id}", taskHandler.Delete)
		r.Get("/my", taskHandler.GetMy)
		r.Post("/status", taskHandler.GetStatus)
	})

	// Notification
	router.Route("/v1/notifications", func(r chi.Router) {
		r.Use(SessionAuthMiddleware(grpcClient))
		r.Get("/my", notificationHandler.GetUserNotifications)
	})

	return router
}

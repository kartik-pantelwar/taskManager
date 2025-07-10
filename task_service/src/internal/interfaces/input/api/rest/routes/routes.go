package routes

import (
	"context"
	"fmt"
	"net/http"
	taskhandler "task_service/src/internal/interfaces/input/api/rest/handler"
	pb "task_service/src/internal/interfaces/input/grpc/generated"

	"github.com/go-chi/chi/v5"
)

func SessionAuthMiddleware(grpcClient pb.SessionValidatorClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ^Extract "sess" cookie (refresh token/session ID)
			cookie, err := r.Cookie("sess")
			if err != nil {
				http.Error(w, "session cookie is missing", http.StatusUnauthorized)
				return
			}
			sessionID := cookie.Value
			fmt.Println("session id : ", sessionID)

			// *gRPC call to user service for session validation
			resp, err := grpcClient.ValidateSession(context.Background(), &pb.ValidateSessionRequest{
				SessionId: sessionID,
			})
			if err != nil || !resp.Valid {
				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}

			// &adding user_id to context for other operations
			ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func InitRoutes(taskHandler *taskhandler.TaskHandler, grpcClient pb.SessionValidatorClient) http.Handler {
	router := chi.NewRouter()
	router.Route("/task", func(r chi.Router) {
		r.Use(SessionAuthMiddleware(grpcClient))
		r.Post("/create", taskHandler.Create)
		r.Put("/update", taskHandler.Update)
		r.Get("/get/all", taskHandler.GetAll)
	})

	return router
}

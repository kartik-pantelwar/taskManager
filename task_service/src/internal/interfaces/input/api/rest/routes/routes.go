package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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
			user_id, err := strconv.Atoi(resp.UserId)
			if err != nil {
				fmt.Printf("Invalid Format: %v", err)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", user_id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func InitRoutes(taskHandler *taskhandler.TaskHandler, grpcClient pb.SessionValidatorClient) http.Handler {
	router := chi.NewRouter()
	router.Route("/v1/tasks", func(r chi.Router) {
		r.Use(SessionAuthMiddleware(grpcClient))
		r.Post("/create", taskHandler.CreateHandler)
		r.Put("/update", taskHandler.UpdateHandler)
		r.Get("/", taskHandler.GetMyHandler)
		r.Delete("/delete/{task-id}",taskHandler.DeleteHandler)
	})
	//task assign krne se pehle ek /status api create krenge jo check kregi ki user available hai, ya nahi, agar available nhi hai, to notify kr dega
	return router
}

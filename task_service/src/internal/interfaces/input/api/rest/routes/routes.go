package routes

import (
	"net/http"

	taskhandler "task_service/src/internal/interfaces/input/api/rest/handler"
	"task_service/src/internal/interfaces/input/api/rest/middleware"
	pb "task_service/src/internal/interfaces/input/grpc/generated/generated"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(taskHandler *taskhandler.TaskHandler, grpcClient pb.SessionValidatorClient) http.Handler {
	router := chi.NewRouter()

	router.Route("/v1/tasks", func(r chi.Router) {
		r.Use(middleware.SessionAuthMiddleware(grpcClient))
		r.Post("/create", taskHandler.Create)
		r.Put("/update", taskHandler.Update)
		r.Delete("/delete/{id}", taskHandler.Delete)
		r.Get("/my", taskHandler.GetMy)
		r.Post("/status", taskHandler.GetStatus)
	})

	return router
}

package routes

import (
	userhandler "user_service/src/internal/interfaces/input/api/rest/handler"
	"user_service/src/internal/interfaces/input/api/rest/middleware"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(
	userHandler *userhandler.UserHandler) http.Handler {
	router := chi.NewRouter()

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)
	})

	router.Route("/users", func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Get("/profile", userHandler.Profile)
		r.Post("/logout", userHandler.LogOut)
		r.Get("/",userHandler.GetAll)
	})

	return router
}

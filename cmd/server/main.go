package main

import (
	"TaskManager/internal/adaptors/persistance"
	"TaskManager/internal/config"
	userhandler "TaskManager/internal/interfaces/input/api/rest/handler/usehandler"
	"TaskManager/internal/interfaces/input/api/rest/routes"
	user "TaskManager/internal/usecase"
	"TaskManager/pkg/migrate"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	database, err := persistance.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}
	fmt.Println("Connected to database")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory %v", err)
	}

	migrate := migrate.NewMigrate(
		database.GetDB(),
		cwd+"/migrations")

	err = migrate.RunMigrations()
	if err != nil {
		log.Fatalf("failed to run migrations %v", err)
	}

	userRepo := persistance.NewUserRepo(database)
	sessionRepo := persistance.NewSessionRepo(database)
	userService := user.NewUserService(userRepo, sessionRepo)
	userHandler := userhandler.NewUserHandler(userService)

	router := routes.InitRoutes(&userHandler)

	configP, err := config.LoadConfig()
	if err != nil {
		panic("Unable to use port")
	}
	err = http.ListenAndServe(fmt.Sprintf(":%s", configP.APP_PORT), router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

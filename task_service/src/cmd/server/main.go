package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"task_service/src/internal/adaptors/persistance"
	"task_service/src/internal/config"
	taskhandler "task_service/src/internal/interfaces/input/api/rest/handler"
	"task_service/src/internal/interfaces/input/api/rest/routes"
	clientpkg "task_service/src/internal/interfaces/input/grpc/client"
	task "task_service/src/internal/usecase"
	"task_service/src/pkg/migrate"
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
		cwd+"/src/migrations")

	err = migrate.RunMigrations()
	if err != nil {
		log.Fatalf("failed to run migrations %v", err)
	}

	configP, err := config.LoadConfig()
	if err != nil {
		panic("Unable to use port")
	}

	grpcClient, err := clientpkg.NewSessionValidatorClient(fmt.Sprintf("localhost:%s", configP.GRPC_PORT))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}

	taskRepo := persistance.NewTaskRepo(database)
	taskService := task.NewTaskService(taskRepo)
	taskHandler := taskhandler.NewTaskHandler(taskService)

	router := routes.InitRoutes(&taskHandler, grpcClient)

	err = http.ListenAndServe(fmt.Sprintf(":%s", configP.APP_PORT), router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

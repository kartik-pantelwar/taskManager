package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"task_service/src/internal/adaptors/persistance"
	redisclient "task_service/src/internal/adaptors/redis"
	client "task_service/src/internal/adaptors/user_grpc_client"
	"task_service/src/internal/config"
	"task_service/src/internal/adaptors/redis/notification"
	taskhandler "task_service/src/internal/interfaces/input/api/rest/handler"
	"task_service/src/internal/interfaces/input/api/rest/routes"
	task "task_service/src/internal/usecase"
	"task_service/src/pkg/migrate"
)

func main() {
	//database setup
	database, err := persistance.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}
	fmt.Println("Connected to database")

	//migration setup
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

	//config setup
	configP, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var notificationService *notification.NotificationService
	redisClient, err := redisclient.NewRedisClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Notifications will be disabled.", err)
		notificationService = nil // Disable notifications if Redis fails
	} else {
		defer redisClient.Close()
		notificationService = notification.NewNotificationService(redisClient.GetClient())
		fmt.Println("Connected to Redis")
	}

	//gRPC client setup
	grpcClient, err := client.NewSessionValidatorClient(fmt.Sprintf("localhost:%s", configP.GRPC_PORT))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}

	taskRepo := persistance.NewTaskRepo(database)
	taskService := task.NewTaskService(taskRepo, notificationService, grpcClient) //added notificationService and grpcClient
	taskHandler := taskhandler.NewTaskHandler(taskService)

	//notification handler - now calls notification service via HTTP
	notificationServiceURL := fmt.Sprintf("http://localhost:%s", configP.NOTIFICATION_PORT) // notification service URL
	notificationHandler := taskhandler.NewNotificationHandler(notificationServiceURL)

	router := routes.InitRoutes(&taskHandler, notificationHandler, grpcClient)

	// server starting
	fmt.Printf("Starting server on port %s\n", configP.APP_PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%s", configP.APP_PORT), router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

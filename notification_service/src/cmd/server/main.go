package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"notificationservice/src/internal/adaptors/redis"
	"notificationservice/src/internal/config"
	"notificationservice/src/internal/interfaces/http/handler"
	"notificationservice/src/internal/interfaces/http/routes"
	"notificationservice/src/internal/interfaces/subscriber"
	"notificationservice/src/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Redis
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("Connected to Redis")

	// Use Redis directly instead of repository
	notificationUseCase := usecase.NewNotificationUseCase(redisClient)
	notificationHandler := handler.NewNotificationHandler(notificationUseCase)
	eventSubscriber := subscriber.NewEventSubscriber(redisClient, notificationUseCase)

	// Start event subscriber in a goroutine
	channel := "task_events" // ^Channel Name of Subscribed Channel
	go eventSubscriber.StartListening(context.Background(), channel)

	// Initialize HTTP routes
	router := routes.InitRoutes(notificationHandler)

	// Start HTTP server
	log.Printf("Notification service HTTP server starting on port %s", cfg.APP_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.APP_PORT), router))
}

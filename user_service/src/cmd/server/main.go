package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"user_service/src/internal/adaptors/persistance"
	"user_service/src/internal/config"
	pb "user_service/src/internal/interfaces/grpc/generated/generated"
	grpcserver "user_service/src/internal/interfaces/grpc/server"
	userhandler "user_service/src/internal/interfaces/input/api/rest/handler"
	"user_service/src/internal/interfaces/input/api/rest/routes"
	user "user_service/src/internal/usecase"
	"user_service/src/pkg/migrate"

	"google.golang.org/grpc"
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

	userRepo := persistance.NewUserRepo(database)
	sessionRepo := persistance.NewSessionRepo(database)
	userService := user.NewUserService(userRepo, sessionRepo)
	userHandler := userhandler.NewUserHandler(userService)

	router := routes.InitRoutes(&userHandler)

	configP, err := config.LoadConfig()
	if err != nil {
		panic("Unable to use port")
	}

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configP.GRPC_PORT))
		if err != nil {
			log.Fatalf("failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		sessionValidatorServer := grpcserver.NewSessionValidatorServer(sessionRepo)
		pb.RegisterSessionValidatorServer(grpcServer, sessionValidatorServer)

		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	log.Printf("HTTP server listening on port %s", configP.APP_PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%s", configP.APP_PORT), router)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net"
	"os"

	"chat-service/internal/handler"
	"chat-service/internal/repository"
	"chat-service/internal/service"
	"chat-service/pkg/centrifugo"
	"chat-service/pkg/database"

	"google.golang.org/grpc"

	"chat-service/internal/proto"
)

func main() {
	// Initialize database connection
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	err = database.AutoMigrate(db)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize Centrifugo client
	centrifugoClient, err := centrifugo.NewCentrifugoClient()
	if err != nil {
		log.Fatalf("Failed to connect to Centrifugo: %v", err)
	}
	defer centrifugoClient.Close()

	// Initialize repository
	chatRepo := repository.NewChatRepository(db)

	// Initialize service
	chatService := service.NewChatService(chatRepo, centrifugoClient)

	// Initialize gRPC handler
	chatHandler := handler.NewChatHandler(chatService)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	// Register chat service
	proto.RegisterChatServiceServer(grpcServer, chatHandler)

	// Start gRPC server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Chat Service gRPC server listening on port %s", port)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// unaryInterceptor is a simple logging interceptor
func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("gRPC method called: %s", info.FullMethod)
	return handler(ctx, req)
}
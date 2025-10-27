package main

import (
	"context"
	"log"
	"net"
	"os"

	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/database"
	"google.golang.org/grpc"

	"auth-service/internal/proto"
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

	// Initialize repository
	authRepo := repository.NewAuthRepository(db)

	// Initialize service
	authService := service.NewAuthService(authRepo)

	// Initialize gRPC handler
	authHandler := handler.NewAuthHandler(authService)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	// Register auth service
	proto.RegisterAuthServiceServer(grpcServer, authHandler)

	// Start gRPC server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Auth Service gRPC server listening on port %s", port)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// unaryInterceptor is a simple logging interceptor
func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("gRPC method called: %s", info.FullMethod)
	return handler(ctx, req)
}
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/pkg/clients"
)

func main() {
	// Initialize gRPC clients
	authClient, err := clients.NewAuthClient()
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authClient.Close()

	// Initialize HTTP handler
	httpHandler := handler.NewHTTPHandler(authClient)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Chat App API Gateway",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Public routes
	public := app.Group("/api/v1")
	public.Get("/health", httpHandler.HealthCheck)
	public.Post("/auth/register", httpHandler.Register)
	public.Post("/auth/login", httpHandler.Login)
	public.Post("/auth/validate", httpHandler.ValidateToken)

	// Protected routes (require authentication)
	protected := app.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authClient))

	// Auth routes
	protected.Get("/auth/profile", httpHandler.GetProfile)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("API Gateway started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")
}
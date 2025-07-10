package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rancher-manager/internal/authservice/grpc"
	"rancher-manager/internal/authservice/handler"
	"rancher-manager/internal/authservice/model"
	"rancher-manager/internal/authservice/repository"
	"rancher-manager/internal/authservice/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=auth_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "authservice"})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/profile", authHandler.AuthMiddleware(), authHandler.GetProfile)
		auth.PUT("/profile", authHandler.AuthMiddleware(), authHandler.UpdateProfile)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Start HTTP server in goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8081"
		}

		log.Printf("AuthService HTTP starting on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	// Start gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	log.Printf("AuthService gRPC starting on port %s", grpcPort)
	if err := grpc.StartGRPCServer(authService, grpcPort); err != nil {
		log.Fatal("Failed to start gRPC server:", err)
	}
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rancher-manager/internal/inventoryservice/grpc"
	"rancher-manager/internal/inventoryservice/handler"
	"rancher-manager/internal/inventoryservice/model"
	"rancher-manager/internal/inventoryservice/repository"
	"rancher-manager/internal/inventoryservice/service"
	"rancher-manager/kafka"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=inventory_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&model.Inventory{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Kafka publisher
	kafkaBrokers := []string{"localhost:9092"}
	publisher, err := kafka.NewPublisher(kafkaBrokers)
	if err != nil {
		log.Printf("Warning: Failed to connect to Kafka: %v", err)
		publisher = nil
	}

	// Initialize gRPC clients
	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authServiceAddr == "" {
		authServiceAddr = "localhost:50051"
	}

	itemServiceAddr := os.Getenv("ITEM_SERVICE_ADDR")
	if itemServiceAddr == "" {
		itemServiceAddr = "localhost:50052"
	}

	authClient, err := grpc.NewAuthClient(authServiceAddr)
	if err != nil {
		log.Fatal("Failed to connect to auth service:", err)
	}

	itemClient, err := grpc.NewItemClient(itemServiceAddr)
	if err != nil {
		log.Fatal("Failed to connect to item service:", err)
	}

	// Initialize layers
	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepo, authClient, itemClient, publisher)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

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
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "inventoryservice"})
	})

	// Inventory routes (all require authentication)
	inventory := r.Group("/inventory")
	inventory.Use(inventoryHandler.AuthMiddleware())
	{
		inventory.POST("/stock/:item_id", inventoryHandler.UpdateStock)
		inventory.DELETE("/item/:item_id", inventoryHandler.DeleteItem)
		inventory.GET("/stock/:item_id", inventoryHandler.GetStock)
		inventory.GET("/items", inventoryHandler.GetAllItems)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("InventoryService starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"rancher-manager/internal/itemservice/grpc"
	"rancher-manager/internal/itemservice/handler"
	"rancher-manager/internal/itemservice/model"
	"rancher-manager/internal/itemservice/repository"
	"rancher-manager/internal/itemservice/service"
	"rancher-manager/kafka"
)

// KafkaEventHandler implements the EventHandler interface for ItemService
type KafkaEventHandler struct {
	itemService *service.ItemService
}

func (h *KafkaEventHandler) HandleStockUpdate(event *kafka.StockUpdateEvent) error {
	log.Printf("Received stock update event for item %s: new stock %d", event.ItemID, event.NewStock)

	// Update stock in ItemService
	updateReq := &model.UpdateItemRequest{
		Stock: event.NewStock,
	}

	_, err := h.itemService.UpdateItem(event.ItemID, updateReq, event.UserID)
	if err != nil {
		log.Printf("Failed to update stock for item %s: %v", event.ItemID, err)
		return err
	}

	log.Printf("Successfully updated stock for item %s to %d", event.ItemID, event.NewStock)
	return nil
}

func (h *KafkaEventHandler) HandleItemDelete(event *kafka.ItemDeleteEvent) error {
	log.Printf("Received item delete event for item %s", event.ItemID)

	// Delete item in ItemService
	err := h.itemService.DeleteItem(event.ItemID, event.UserID)
	if err != nil {
		log.Printf("Failed to delete item %s: %v", event.ItemID, err)
		return err
	}

	log.Printf("Successfully deleted item %s", event.ItemID)
	return nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// MongoDB connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB")

	// Initialize database
	db := client.Database("item_db")

	// Initialize gRPC auth client
	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authServiceAddr == "" {
		authServiceAddr = "localhost:50051"
	}

	authClient, err := grpc.NewAuthClient(authServiceAddr)
	if err != nil {
		log.Fatal("Failed to connect to auth service:", err)
	}

	// Initialize layers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo, authClient)
	itemHandler := handler.NewItemHandler(itemService)

	// Initialize Kafka consumer
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	kafkaHandler := &KafkaEventHandler{itemService: itemService}
	consumer, err := kafka.NewConsumer([]string{kafkaBrokers}, "item-service-group", kafkaHandler)
	if err != nil {
		log.Printf("Warning: Failed to connect to Kafka: %v", err)
	} else {
		// Start Kafka consumer in goroutine
		go func() {
			topics := []string{"stock_updates", "item_deletes"}
			if err := consumer.Start(context.Background(), topics); err != nil {
				log.Printf("Kafka consumer error: %v", err)
			}
		}()
		log.Println("Kafka consumer started")
	}

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
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "itemservice"})
	})

	// Item routes (all require authentication)
	items := r.Group("/items")
	items.Use(itemHandler.AuthMiddleware())
	{
		items.POST("/", itemHandler.CreateItem)
		items.GET("/", itemHandler.GetAllItems)
		items.GET("/:id", itemHandler.GetItem)
		items.PUT("/:id", itemHandler.UpdateItem)
		items.DELETE("/:id", itemHandler.DeleteItem)
	}

	// Start HTTP server in goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8082"
		}

		log.Printf("ItemService HTTP starting on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	// Start gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	log.Printf("ItemService gRPC starting on port %s", grpcPort)
	if err := grpc.StartGRPCServer(itemService, grpcPort); err != nil {
		log.Fatal("Failed to start gRPC server:", err)
	}
}

package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Rancher Manager API Gateway",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"service":   "api-gateway",
			"version":   "1.0.0",
			"framework": "fiber",
		})
	})

	// Service URLs
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8081"
	}

	itemServiceURL := os.Getenv("ITEM_SERVICE_URL")
	if itemServiceURL == "" {
		itemServiceURL = "http://localhost:8082"
	}

	inventoryServiceURL := os.Getenv("INVENTORY_SERVICE_URL")
	if inventoryServiceURL == "" {
		inventoryServiceURL = "http://localhost:8083"
	}

	// Auth Service Proxy
	app.All("/auth/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		targetURL := authServiceURL + "/" + path

		return proxyRequest(c, targetURL)
	})

	// Item Service Proxy
	app.All("/items/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		targetURL := itemServiceURL + "/" + path

		return proxyRequest(c, targetURL)
	})

	// Inventory Service Proxy
	app.All("/inventory/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		targetURL := inventoryServiceURL + "/" + path

		return proxyRequest(c, targetURL)
	})

	// API Documentation
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "Rancher Manager API Gateway",
			"version":   "1.0.0",
			"framework": "fiber",
			"services": fiber.Map{
				"auth": fiber.Map{
					"base_url":    "/auth",
					"service_url": authServiceURL,
					"endpoints": []string{
						"POST /auth/register",
						"POST /auth/login",
						"POST /auth/logout",
						"GET /auth/profile",
						"PUT /auth/profile",
						"POST /auth/refresh",
						"GET /auth/health",
					},
				},
				"items": fiber.Map{
					"base_url":    "/items",
					"service_url": itemServiceURL,
					"endpoints": []string{
						"POST /items/",
						"GET /items/",
						"GET /items/:id",
						"PUT /items/:id",
						"DELETE /items/:id",
						"GET /items/health",
					},
				},
				"inventory": fiber.Map{
					"base_url":    "/inventory",
					"service_url": inventoryServiceURL,
					"endpoints": []string{
						"POST /inventory/stock/:item_id",
						"GET /inventory/stock/:item_id",
						"GET /inventory/items",
						"DELETE /inventory/item/:item_id",
						"GET /inventory/health",
					},
				},
			},
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	log.Printf("Auth Service: %s", authServiceURL)
	log.Printf("Item Service: %s", itemServiceURL)
	log.Printf("Inventory Service: %s", inventoryServiceURL)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}

// proxyRequest forwards the request to the target service
func proxyRequest(c *fiber.Ctx, targetURL string) error {
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(c.Method(), targetURL, bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create request",
		})
	}

	// Copy headers
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set forwarded headers
	req.Header.Set("X-Forwarded-Host", c.Hostname())
	req.Header.Set("X-Forwarded-Proto", c.Protocol())
	req.Header.Set("X-Forwarded-For", c.IP())

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to forward request",
		})
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read response",
		})
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Return response
	return c.Status(resp.StatusCode).Send(body)
}

package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"rancher-manager/internal/inventoryservice/grpc"
	"rancher-manager/internal/inventoryservice/model"
	"rancher-manager/internal/inventoryservice/service"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
	authClient       *grpc.AuthClient
}

func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		authClient:       inventoryService.GetAuthClient(),
	}
}

// UpdateStock godoc
// @Summary Update item stock
// @Description Update stock for a specific item
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id path string true "Item ID"
// @Param stock body model.UpdateStockRequest true "Stock update data"
// @Success 200 {object} model.StockResponse
// @Failure 400 {object} model.StockResponse
// @Failure 401 {object} model.StockResponse
// @Failure 404 {object} model.StockResponse
// @Router /inventory/stock/{item_id} [post]
func (h *InventoryHandler) UpdateStock(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.StockResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, model.StockResponse{
			Message: "Item ID is required",
			Success: false,
		})
		return
	}

	var req model.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.StockResponse{
			Message: "Invalid request data: " + err.Error(),
			Success: false,
		})
		return
	}

	inventory, err := h.inventoryService.UpdateStock(itemID, req.NewStock, userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, model.StockResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, model.StockResponse{
		Message: "Stock updated successfully",
		Success: true,
		Data:    inventory,
	})
}

// DeleteItem godoc
// @Summary Delete item from inventory
// @Description Delete an item from inventory and item service
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id path string true "Item ID"
// @Success 200 {object} model.StockResponse
// @Failure 400 {object} model.StockResponse
// @Failure 401 {object} model.StockResponse
// @Failure 404 {object} model.StockResponse
// @Router /inventory/item/{item_id} [delete]
func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.StockResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, model.StockResponse{
			Message: "Item ID is required",
			Success: false,
		})
		return
	}

	err := h.inventoryService.DeleteItem(itemID, userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, model.StockResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, model.StockResponse{
		Message: "Item deleted successfully",
		Success: true,
	})
}

// GetStock godoc
// @Summary Get item stock
// @Description Get stock information for a specific item
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id path string true "Item ID"
// @Success 200 {object} model.StockResponse
// @Failure 400 {object} model.StockResponse
// @Failure 401 {object} model.StockResponse
// @Failure 404 {object} model.StockResponse
// @Router /inventory/stock/{item_id} [get]
func (h *InventoryHandler) GetStock(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.StockResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, model.StockResponse{
			Message: "Item ID is required",
			Success: false,
		})
		return
	}

	inventory, err := h.inventoryService.GetStock(itemID, userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, model.StockResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, model.StockResponse{
		Message: "Stock retrieved successfully",
		Success: true,
		Data:    inventory,
	})
}

// GetAllItems godoc
// @Summary Get all inventory items
// @Description Get all inventory items with stock information
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.ItemsResponse
// @Failure 401 {object} model.ItemsResponse
// @Router /inventory/items [get]
func (h *InventoryHandler) GetAllItems(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ItemsResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	inventories, err := h.inventoryService.GetAllItems(userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		}
		c.JSON(status, model.ItemsResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, model.ItemsResponse{
		Message: "Inventory items retrieved successfully",
		Success: true,
		Data:    inventories,
	})
}

// AuthMiddleware validates JWT token via gRPC and sets user context
func (h *InventoryHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.StockResponse{
				Message: "Authorization header required",
				Success: false,
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.StockResponse{
				Message: "Invalid authorization header format",
				Success: false,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token via gRPC
		response, err := h.authClient.ValidateToken(tokenString)
		if err != nil || !response.Valid {
			c.JSON(http.StatusUnauthorized, model.StockResponse{
				Message: "Invalid or expired token",
				Success: false,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", response.UserId)
		c.Set("username", response.Username)
		c.Set("role", response.Role)

		c.Next()
	}
}

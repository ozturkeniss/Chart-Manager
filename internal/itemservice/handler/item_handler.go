package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"rancher-manager/internal/itemservice/model"
	"rancher-manager/internal/itemservice/service"
)

type AuthClientInterface interface {
	ValidateToken(token string) (interface{}, error)
	GetUser(userID uint32) (interface{}, error)
}

type ItemHandler struct {
	itemService *service.ItemService
	authClient  AuthClientInterface
}

func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
		authClient:  itemService.GetAuthClient(),
	}
}

// CreateItem godoc
// @Summary Create a new item
// @Description Create a new item with authentication
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body model.CreateItemRequest true "Item data"
// @Success 201 {object} model.ItemResponse
// @Failure 400 {object} model.ItemResponse
// @Failure 401 {object} model.ItemResponse
// @Router /items/ [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ItemResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	var req model.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ItemResponse{
			Message: "Invalid request data: " + err.Error(),
			Success: false,
		})
		return
	}

	item, err := h.itemService.CreateItem(&req, userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		}
		c.JSON(status, model.ItemResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusCreated, model.ItemResponse{
		Message: "Item created successfully",
		Success: true,
		Data:    item,
	})
}

// GetItem godoc
// @Summary Get item by ID
// @Description Get a specific item by ID with authentication
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Item ID"
// @Success 200 {object} model.ItemResponse
// @Failure 400 {object} model.ItemResponse
// @Failure 401 {object} model.ItemResponse
// @Failure 404 {object} model.ItemResponse
// @Router /items/{id} [get]
func (h *ItemHandler) GetItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ItemResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ItemResponse{
			Message: "Item ID is required",
			Success: false,
		})
		return
	}

	item, err := h.itemService.GetItem(id, userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, model.ItemResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, model.ItemResponse{
		Message: "Item retrieved successfully",
		Success: true,
		Data:    item,
	})
}

// GetAllItems godoc
// @Summary Get all items
// @Description Get all items with authentication
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.ItemsResponse
// @Failure 401 {object} model.ItemsResponse
// @Router /items/ [get]
func (h *ItemHandler) GetAllItems(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ItemsResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	items, err := h.itemService.GetAllItems(userID.(uint32))
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
		Message: "Items retrieved successfully",
		Success: true,
		Data:    items,
	})
}

// UpdateItem godoc
// @Summary Update item
// @Description Update an item by ID with authentication
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Item ID"
// @Param item body model.UpdateItemRequest true "Item update data"
// @Success 200 {object} model.ItemResponse
// @Failure 400 {object} model.ItemResponse
// @Failure 401 {object} model.ItemResponse
// @Failure 404 {object} model.ItemResponse
// @Router /items/{id} [put]
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ItemResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ItemResponse{
			Message: "Item ID is required",
			Success: false,
		})
		return
	}

	var req model.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ItemResponse{
			Message: "Invalid request data: " + err.Error(),
			Success: false,
		})
		return
	}

	item, err := h.itemService.UpdateItem(id, &req, userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, model.ItemResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, model.ItemResponse{
		Message: "Item updated successfully",
		Success: true,
		Data:    item,
	})
}

// DeleteItem godoc
// @Summary Delete item
// @Description Delete an item by ID with authentication
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Item ID"
// @Success 200 {object} model.ItemResponse
// @Failure 400 {object} model.ItemResponse
// @Failure 401 {object} model.ItemResponse
// @Failure 404 {object} model.ItemResponse
// @Router /items/{id} [delete]
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ItemResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ItemResponse{
			Message: "Item ID is required",
			Success: false,
		})
		return
	}

	err := h.itemService.DeleteItem(id, userID.(uint32))
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, model.ItemResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, model.ItemResponse{
		Message: "Item deleted successfully",
		Success: true,
	})
}

// AuthMiddleware validates JWT token via gRPC and sets user context
func (h *ItemHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ItemResponse{
				Message: "Authorization header required",
				Success: false,
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.ItemResponse{
				Message: "Invalid authorization header format",
				Success: false,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token via gRPC
		response, err := h.authClient.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ItemResponse{
				Message: "Invalid or expired token",
				Success: false,
			})
			c.Abort()
			return
		}

		// Type assertion for response - use interface{} to avoid import cycle
		if validateResponse, ok := response.(map[string]interface{}); ok {
			if valid, exists := validateResponse["valid"].(bool); exists && valid {
				// Set user information in context
				if userID, exists := validateResponse["user_id"].(uint32); exists {
					c.Set("user_id", userID)
				}
				if username, exists := validateResponse["username"].(string); exists {
					c.Set("username", username)
				}
				if role, exists := validateResponse["role"].(string); exists {
					c.Set("role", role)
				}
			} else {
				c.JSON(http.StatusUnauthorized, model.ItemResponse{
					Message: "Invalid or expired token",
					Success: false,
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, model.ItemResponse{
				Message: "Invalid or expired token",
				Success: false,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

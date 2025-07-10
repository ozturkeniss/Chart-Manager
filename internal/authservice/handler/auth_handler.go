package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"rancher-manager/internal/authservice/model"
	"rancher-manager/internal/authservice/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, password, first name and last name
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "User registration data"
// @Success 201 {object} model.User
// @Failure 400 {object} model.AuthResponse
// @Failure 409 {object} model.AuthResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.AuthResponse{
			Message: "Invalid request data: " + err.Error(),
			Success: false,
		})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}
		c.JSON(status, model.AuthResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Login user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.AuthResponse
// @Failure 401 {object} model.AuthResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.AuthResponse{
			Message: "Invalid request data: " + err.Error(),
			Success: false,
		})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.AuthResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user (client should discard tokens)
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.AuthResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// by discarding the token. This endpoint is for consistency.
	c.JSON(http.StatusOK, model.AuthResponse{
		Message: "Successfully logged out",
		Success: true,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User
// @Failure 401 {object} model.AuthResponse
// @Failure 404 {object} model.AuthResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.AuthResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	user, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, model.AuthResponse{
			Message: "User not found",
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body model.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} model.User
// @Failure 400 {object} model.AuthResponse
// @Failure 401 {object} model.AuthResponse
// @Failure 404 {object} model.AuthResponse
// @Failure 409 {object} model.AuthResponse
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.AuthResponse{
			Message: "User not authenticated",
			Success: false,
		})
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.AuthResponse{
			Message: "Invalid request data: " + err.Error(),
			Success: false,
		})
		return
	}

	user, err := h.authService.UpdateProfile(userID.(uint), &req)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}
		c.JSON(status, model.AuthResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.AuthResponse
// @Failure 401 {object} model.AuthResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.AuthResponse{
			Message: "Invalid request data: " + err.Error(),
			Success: false,
		})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.AuthResponse{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthMiddleware validates JWT token and sets user context
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.AuthResponse{
				Message: "Authorization header required",
				Success: false,
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.AuthResponse{
				Message: "Invalid authorization header format",
				Success: false,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.AuthResponse{
				Message: "Invalid or expired token",
				Success: false,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

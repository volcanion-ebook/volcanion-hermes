package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/volcanion/volcanion-hermes/internal/middleware"
	"github.com/volcanion/volcanion-hermes/internal/models"
	"github.com/volcanion/volcanion-hermes/internal/services"
)

// AuthHandler xử lý các HTTP requests liên quan đến authentication
// Bao gồm register, login, profile, và token refresh
type AuthHandler struct {
	userService   *services.UserService      // Service xử lý user business logic
	jwtMiddleware *middleware.JWTMiddleware   // Middleware xử lý JWT operations
}

// NewAuthHandler tạo instance mới của AuthHandler
func NewAuthHandler(userService *services.UserService, jwtMiddleware *middleware.JWTMiddleware) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		jwtMiddleware: jwtMiddleware,
	}
}

// Register xử lý đăng ký tài khoản mới
// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	
	// Bind và validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Tạo user mới qua service layer
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		// Trả về Conflict nếu username/email đã tồn tại
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	// Tạo JWT token cho user mới
	token, err := h.jwtMiddleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	// Trả về user info và token
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data: gin.H{
			"user":  user,
			"token": token,
		},
	})
}

// Login xử lý đăng nhập
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	
	// Bind và validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Tìm user theo username
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil {
		// Không tiết lộ chi tiết lỗi để tránh username enumeration attack
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
			Error:   "invalid_credentials",
		})
		return
	}

	// Verify password
	if err := h.userService.ValidatePassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
			Error:   "invalid_credentials",
		})
		return
	}

	// Kiểm tra account có active không
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User account is inactive",
			Error:   "inactive_account",
		})
		return
	}

	// Tạo JWT token cho user
	token, err := h.jwtMiddleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: gin.H{
			"user":  user,
			"token": token,
		},
	})
}

// GetProfile lấy thông tin profile của user đã đăng nhập
// GET /auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Lấy thông tin user từ context thông qua JWT middleware helper
	user, err := h.jwtMiddleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

// RefreshToken tạo mới JWT token với thông tin user cập nhật
// POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Lấy thông tin user từ token hiện tại
	user, err := h.jwtMiddleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   err.Error(),
		})
		return
	}

	// Convert user ID từ string sang ObjectID
	userID, err := primitive.ObjectIDFromHex(user.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
		return
	}

	// Lấy thông tin user mới nhất từ database để đảm bảo token chứa data hiện tại
	freshUser, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	// Tạo token mới với thông tin user cập nhật
	newToken, err := h.jwtMiddleware.GenerateToken(freshUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: gin.H{
			"token": newToken,
		},
	})
}

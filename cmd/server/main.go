package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/volcanion/volcanion-hermes/internal/config"
	"github.com/volcanion/volcanion-hermes/internal/database"
	"github.com/volcanion/volcanion-hermes/internal/handlers"
	"github.com/volcanion/volcanion-hermes/internal/middleware"
	"github.com/volcanion/volcanion-hermes/internal/models"
	"github.com/volcanion/volcanion-hermes/internal/services"
	"github.com/volcanion/volcanion-hermes/internal/storage"
)

// main là entry point của ứng dụng Volcanion Hermes
// Khởi tạo tất cả dependencies và start HTTP server
func main() {
	// Load configuration từ file .env và environment variables
	cfg := config.Load()

	// Set Gin mode theo config (debug/release/test)
	gin.SetMode(cfg.Server.GinMode)

	// Khởi tạo kết nối MongoDB với connection pooling và indexes
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Khởi tạo MinIO storage service cho quản lý files
	minioStorage, err := storage.NewMinIOStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	// Khởi tạo business logic services
	userService := services.NewUserService(db)
	ebookService := services.NewEbookService(db, minioStorage)

	// Khởi tạo JWT middleware với secret key và expiration config
	jwtMiddleware := middleware.NewJWTMiddleware(cfg)

	// Khởi tạo HTTP handlers với dependency injection
	authHandler := handlers.NewAuthHandler(userService, jwtMiddleware)
	ebookHandler := handlers.NewEbookHandler(ebookService, jwtMiddleware)

	// Setup router với middleware và routes
	router := setupRouter(authHandler, ebookHandler, jwtMiddleware)

	// Start HTTP server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter cấu hình Gin router với middleware và routes
func setupRouter(authHandler *handlers.AuthHandler, ebookHandler *handlers.EbookHandler, jwtMiddleware *middleware.JWTMiddleware) *gin.Engine {
	router := gin.Default()

	// CORS middleware để cho phép cross-origin requests
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check endpoint để monitoring
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "Server is running",
			Data: gin.H{
				"status": "healthy",
			},
		})
	})

	// API routes với versioning
	api := router.Group("/api/v1")

	// Authentication routes - không yêu cầu auth
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/profile", jwtMiddleware.AuthRequired(), authHandler.GetProfile)
		auth.POST("/refresh", jwtMiddleware.AuthRequired(), authHandler.RefreshToken)
	}

	// Ebook routes với phân quyền chi tiết
	ebooks := api.Group("/ebooks")
	{
		// Public routes - không cần authentication
		ebooks.GET("", ebookHandler.ListEbooks)
		ebooks.GET("/search", ebookHandler.SearchEbooks)
		ebooks.GET("/:id", ebookHandler.GetEbook)
		
		// Download cần authentication để tracking
		ebooks.GET("/:id/download", jwtMiddleware.AuthRequired(), ebookHandler.DownloadEbook)

		// Protected routes - yêu cầu authentication
		protected := ebooks.Group("")
		protected.Use(jwtMiddleware.AuthRequired())
		{
			// Admin và editor có thể thêm, sửa, xóa ebook
			adminEditor := protected.Group("")
			adminEditor.Use(jwtMiddleware.RequireRoles("admin", "editor"))
			{
				adminEditor.POST("", ebookHandler.CreateEbook)
				adminEditor.PUT("/:id", ebookHandler.UpdateEbook)
				adminEditor.DELETE("/:id", ebookHandler.DeleteEbook)
				adminEditor.POST("/:id/upload", ebookHandler.UploadEbookFile)
				adminEditor.POST("/:id/cover", ebookHandler.UploadCoverImage)
			}
		}
	}

	return router
}

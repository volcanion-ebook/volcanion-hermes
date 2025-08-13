package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/volcanion/volcanion-hermes/internal/middleware"
	"github.com/volcanion/volcanion-hermes/internal/models"
	"github.com/volcanion/volcanion-hermes/internal/services"
)

// EbookHandler xử lý các HTTP request liên quan đến ebook
type EbookHandler struct {
	ebookService  *services.EbookService     // Service xử lý business logic ebook
	jwtMiddleware *middleware.JWTMiddleware  // JWT middleware để lấy user info
}

// NewEbookHandler tạo instance mới của EbookHandler
func NewEbookHandler(ebookService *services.EbookService, jwtMiddleware *middleware.JWTMiddleware) *EbookHandler {
	return &EbookHandler{
		ebookService:  ebookService,
		jwtMiddleware: jwtMiddleware,
	}
}

// CreateEbook tạo mới một ebook
// POST /ebooks
func (h *EbookHandler) CreateEbook(c *gin.Context) {
	// Lấy thông tin user từ JWT token
	user, err := h.jwtMiddleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   err.Error(),
		})
		return
	}

	// Bind request body vào struct
	var req models.CreateEbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
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

	// Gọi service để tạo ebook
	ebook, err := h.ebookService.CreateEbook(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create ebook",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Ebook created successfully",
		Data:    ebook,
	})
}

// GetEbook lấy thông tin ebook theo ID
// GET /ebooks/:id
func (h *EbookHandler) GetEbook(c *gin.Context) {
	// Lấy ID từ URL parameter
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ebook ID",
			Error:   err.Error(),
		})
		return
	}

	// Tìm ebook theo ID
	ebook, err := h.ebookService.GetEbookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Ebook not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ebook retrieved successfully",
		Data:    ebook,
	})
}

// UpdateEbook cập nhật thông tin ebook
// PUT /ebooks/:id
func (h *EbookHandler) UpdateEbook(c *gin.Context) {
	// Lấy ID từ URL parameter
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ebook ID",
			Error:   err.Error(),
		})
		return
	}

	// Bind request body
	var req models.UpdateEbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Cập nhật ebook thông qua service
	ebook, err := h.ebookService.UpdateEbook(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update ebook",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ebook updated successfully",
		Data:    ebook,
	})
}

// DeleteEbook xóa ebook theo ID
// DELETE /ebooks/:id
func (h *EbookHandler) DeleteEbook(c *gin.Context) {
	// Lấy ID từ URL parameter
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ebook ID",
			Error:   err.Error(),
		})
		return
	}

	// Xóa ebook - service sẽ xử lý cascade delete cho files
	err = h.ebookService.DeleteEbook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete ebook",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ebook deleted successfully",
	})
}

// ListEbooks lấy danh sách ebook với phân trang và filter
// GET /ebooks?page=1&limit=10&category=fiction&author=john
func (h *EbookHandler) ListEbooks(c *gin.Context) {
	// Parse query parameters cho phân trang
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Giới hạn tối đa 100 items/page để tránh load quá nặng
	}

	// Xây dựng filter từ query parameters
	filter := bson.M{}
	
	// Filter theo category
	if category := c.Query("category"); category != "" {
		filter["category"] = category
	}
	
	// Filter theo author (case-insensitive search)
	if author := c.Query("author"); author != "" {
		filter["author"] = bson.M{"$regex": author, "$options": "i"}
	}

	// Gọi service để lấy danh sách ebook
	response, err := h.ebookService.ListEbooks(page, limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve ebooks",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ebooks retrieved successfully",
		Data:    response,
	})
}

// SearchEbooks tìm kiếm ebook bằng text search
// GET /ebooks/search?q=golang&page=1&limit=10
func (h *EbookHandler) SearchEbooks(c *gin.Context) {
	// Lấy search query từ parameter
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Search query is required",
			Error:   "missing_query",
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Thực hiện tìm kiếm sử dụng MongoDB text index
	response, err := h.ebookService.SearchEbooks(query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to search ebooks",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Search completed successfully",
		Data:    response,
	})
}

// UploadEbookFile upload file ebook (PDF, EPUB, etc.)
// POST /ebooks/:id/upload
func (h *EbookHandler) UploadEbookFile(c *gin.Context) {
	// Lấy ebook ID từ URL parameter
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ebook ID",
			Error:   err.Error(),
		})
		return
	}

	// Lấy file từ form upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "No file uploaded",
			Error:   err.Error(),
		})
		return
	}

	// Upload file thông qua service - service sẽ validate file type và size
	uploadResponse, err := h.ebookService.UploadEbookFile(id, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to upload file",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "File uploaded successfully",
		Data:    uploadResponse,
	})
}

// UploadCoverImage upload ảnh bìa cho ebook
// POST /ebooks/:id/cover
func (h *EbookHandler) UploadCoverImage(c *gin.Context) {
	// Lấy ebook ID từ URL parameter
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ebook ID",
			Error:   err.Error(),
		})
		return
	}

	// Lấy file ảnh từ form upload
	file, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "No cover image uploaded",
			Error:   err.Error(),
		})
		return
	}

	// Upload ảnh bìa - service sẽ validate file là image và resize nếu cần
	uploadResponse, err := h.ebookService.UploadCoverImage(id, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to upload cover image",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Cover image uploaded successfully",
		Data:    uploadResponse,
	})
}

// DownloadEbook tạo presigned URL để download file ebook
// GET /ebooks/:id/download
func (h *EbookHandler) DownloadEbook(c *gin.Context) {
	// Lấy ebook ID từ URL parameter
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ebook ID",
			Error:   err.Error(),
		})
		return
	}

	// Tạo presigned URL cho download file từ MinIO
	downloadURL, err := h.ebookService.GetEbookFile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "File not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Download URL generated successfully",
		Data: gin.H{
			"download_url": downloadURL,
		},
	})
}

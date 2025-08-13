package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/volcanion/volcanion-hermes/internal/database"
	"github.com/volcanion/volcanion-hermes/internal/models"
	"github.com/volcanion/volcanion-hermes/internal/storage"
)

// EbookService xử lý business logic liên quan đến ebook management
// Bao gồm CRUD operations, file upload, và search functionality
type EbookService struct {
	db         *database.Database    // Database connection
	collection *mongo.Collection     // Ebooks collection reference
	storage    *storage.MinIOStorage // File storage service
}

// NewEbookService tạo instance mới của EbookService
func NewEbookService(db *database.Database, storage *storage.MinIOStorage) *EbookService {
	return &EbookService{
		db:         db,
		collection: db.GetCollection("ebooks"),
		storage:    storage,
	}
}

// CreateEbook tạo ebook mới trong hệ thống
// Chỉ tạo metadata, file sẽ được upload riêng sau
func (s *EbookService) CreateEbook(req *models.CreateEbookRequest, userID primitive.ObjectID) (*models.Ebook, error) {
	// Tạo ebook object với metadata từ request
	ebook := &models.Ebook{
		Title:       req.Title,
		Author:      req.Author,
		Publisher:   req.Publisher,
		PublishYear: req.PublishYear,
		ISBN:        req.ISBN,
		Description: req.Description,
		Language:    req.Language,
		Category:    req.Category,
		Tags:        req.Tags,
		CreatedBy:   userID, // Track user tạo ebook
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert ebook vào database
	result, err := s.collection.InsertOne(context.Background(), ebook)
	if err != nil {
		// Kiểm tra duplicate ISBN (nếu có)
		if mongo.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("ebook with this ISBN already exists")
		}
		return nil, fmt.Errorf("failed to create ebook: %w", err)
	}

	// Set ID từ result
	ebook.ID = result.InsertedID.(primitive.ObjectID)
	return ebook, nil
}

// GetEbookByID lấy ebook theo ID
func (s *EbookService) GetEbookByID(id primitive.ObjectID) (*models.Ebook, error) {
	var ebook models.Ebook
	err := s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&ebook)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("ebook not found")
		}
		return nil, fmt.Errorf("failed to get ebook: %w", err)
	}
	return &ebook, nil
}

// UpdateEbook cập nhật ebook với partial updates
// Chỉ update những field được cung cấp (non-nil values)
func (s *EbookService) UpdateEbook(id primitive.ObjectID, req *models.UpdateEbookRequest) (*models.Ebook, error) {
	updates := bson.M{"updated_at": time.Now()}

	// Chỉ thêm field vào updates nếu có giá trị (pointer không nil)
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Author != nil {
		updates["author"] = *req.Author
	}
	if req.Publisher != nil {
		updates["publisher"] = *req.Publisher
	}
	if req.PublishYear != nil {
		updates["publish_year"] = *req.PublishYear
	}
	if req.ISBN != nil {
		updates["isbn"] = *req.ISBN
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Language != nil {
		updates["language"] = *req.Language
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	// Thực hiện update
	_, err := s.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update ebook: %w", err)
	}

	// Trả về ebook đã được update
	return s.GetEbookByID(id)
}

// DeleteEbook xóa ebook và tất cả files liên quan
// Cascade delete: metadata + main file + cover image + page files
func (s *EbookService) DeleteEbook(id primitive.ObjectID) error {
	// Lấy thông tin ebook để biết files cần xóa
	ebook, err := s.GetEbookByID(id)
	if err != nil {
		return err
	}

	// Xóa main file từ storage
	if ebook.FilePath != "" {
		if err := s.storage.DeleteFile(ebook.FilePath); err != nil {
			// Log error nhưng không fail toàn bộ operation
			// File có thể đã bị xóa manual hoặc corrupted
			fmt.Printf("Warning: failed to delete main file %s: %v\n", ebook.FilePath, err)
		}
	}

	// Xóa cover image từ storage
	if ebook.CoverImage != "" {
		if err := s.storage.DeleteFile(ebook.CoverImage); err != nil {
			fmt.Printf("Warning: failed to delete cover image %s: %v\n", ebook.CoverImage, err)
		}
	}

	// Xóa tất cả page files từ storage
	for _, page := range ebook.Pages {
		if page.FilePath != "" {
			if err := s.storage.DeleteFile(page.FilePath); err != nil {
				fmt.Printf("Warning: failed to delete page file %s: %v\n", page.FilePath, err)
			}
		}
	}

	// Cuối cùng xóa metadata từ database
	_, err = s.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete ebook: %w", err)
	}

	return nil
}

func (s *EbookService) ListEbooks(page, limit int, filter bson.M) (*models.EbookListResponse, error) {
	skip := (page - 1) * limit

	// Count total documents
	total, err := s.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count ebooks: %w", err)
	}

	// Find ebooks with pagination
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
	cursor, err := s.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find ebooks: %w", err)
	}
	defer cursor.Close(context.Background())

	var ebooks []models.Ebook
	if err = cursor.All(context.Background(), &ebooks); err != nil {
		return nil, fmt.Errorf("failed to decode ebooks: %w", err)
	}

	return &models.EbookListResponse{
		Ebooks: ebooks,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

func (s *EbookService) SearchEbooks(query string, page, limit int) (*models.EbookListResponse, error) {
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	return s.ListEbooks(page, limit, filter)
}

func (s *EbookService) UploadEbookFile(ebookID primitive.ObjectID, file *multipart.FileHeader) (*models.FileUploadResponse, error) {
	// Validate file type
	if !s.storage.IsValidFileType(file.Filename) {
		return nil, fmt.Errorf("invalid file type")
	}

	// Upload file
	folder := fmt.Sprintf("ebooks/%s", ebookID.Hex())
	uploadResponse, err := s.storage.UploadFile(file, folder)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Update ebook with file information
	updates := bson.M{
		"file_path":   uploadResponse.Filename,
		"file_size":   uploadResponse.Size,
		"file_format": s.getFileFormat(file.Filename),
		"updated_at":  time.Now(),
	}

	_, err = s.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": ebookID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update ebook with file info: %w", err)
	}

	return uploadResponse, nil
}

func (s *EbookService) UploadCoverImage(ebookID primitive.ObjectID, file *multipart.FileHeader) (*models.FileUploadResponse, error) {
	// Upload cover image
	folder := fmt.Sprintf("covers/%s", ebookID.Hex())
	uploadResponse, err := s.storage.UploadFile(file, folder)
	if err != nil {
		return nil, fmt.Errorf("failed to upload cover image: %w", err)
	}

	// Update ebook with cover image path
	_, err = s.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": ebookID},
		bson.M{"$set": bson.M{
			"cover_image": uploadResponse.Filename,
			"updated_at":  time.Now(),
		}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update ebook with cover image: %w", err)
	}

	return uploadResponse, nil
}

func (s *EbookService) AddPages(ebookID primitive.ObjectID, pages []models.EbookPage) error {
	_, err := s.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": ebookID},
		bson.M{
			"$set": bson.M{
				"pages":       pages,
				"total_pages": len(pages),
				"updated_at":  time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add pages to ebook: %w", err)
	}
	return nil
}

func (s *EbookService) GetEbooksByCategory(category string, page, limit int) (*models.EbookListResponse, error) {
	filter := bson.M{"category": category}
	return s.ListEbooks(page, limit, filter)
}

func (s *EbookService) GetEbooksByAuthor(author string, page, limit int) (*models.EbookListResponse, error) {
	filter := bson.M{"author": bson.M{"$regex": author, "$options": "i"}}
	return s.ListEbooks(page, limit, filter)
}

func (s *EbookService) getFileFormat(filename string) string {
	ext := filename[len(filename)-4:]
	switch ext {
	case ".pdf":
		return "PDF"
	case "epub":
		return "EPUB"
	case "mobi":
		return "MOBI"
	default:
		return "Unknown"
	}
}

func (s *EbookService) GetEbookFile(ebookID primitive.ObjectID) (string, error) {
	ebook, err := s.GetEbookByID(ebookID)
	if err != nil {
		return "", err
	}

	if ebook.FilePath == "" {
		return "", fmt.Errorf("no file associated with this ebook")
	}

	// Generate presigned URL for download
	url, err := s.storage.GetPresignedURL(ebook.FilePath, time.Hour*24)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return url, nil
}

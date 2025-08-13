package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/volcanion/volcanion-hermes/internal/config"
	"github.com/volcanion/volcanion-hermes/internal/models"
)

// MinIOStorage wrapper cho MinIO client với các operations cần thiết
type MinIOStorage struct {
	client *minio.Client   // MinIO client để thực hiện operations
	config *config.Config  // App config chứa thông tin kết nối
}

// NewMinIOStorage tạo kết nối mới tới MinIO server
// Khởi tạo bucket nếu chưa tồn tại
func NewMinIOStorage(cfg *config.Config) (*MinIOStorage, error) {
	// Tạo MinIO client với credentials từ config
	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL, // HTTP/HTTPS based on config
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	storage := &MinIOStorage{
		client: client,
		config: cfg,
	}

	// Tạo bucket mặc định nếu chưa tồn tại
	if err := storage.createBucketIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	log.Printf("Successfully connected to MinIO at %s", cfg.MinIO.Endpoint)
	return storage, nil
}

// createBucketIfNotExists kiểm tra và tạo bucket nếu chưa tồn tại
// Bucket là container để lưu trữ objects trong MinIO
func (s *MinIOStorage) createBucketIfNotExists() error {
	ctx := context.Background()
	bucketName := s.config.MinIO.Bucket

	// Kiểm tra bucket đã tồn tại chưa
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	// Tạo bucket mới nếu chưa tồn tại
	if !exists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket: %s", bucketName)
	}

	return nil
}

// UploadFile upload file từ multipart form lên MinIO
// Tạo unique filename và trả về thông tin file đã upload
func (s *MinIOStorage) UploadFile(file *multipart.FileHeader, folder string) (*models.FileUploadResponse, error) {
	// Mở file từ multipart header
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Tạo unique filename với timestamp để tránh trùng lặp
	filename := s.generateUniqueFilename(file.Filename, folder)
	objectName := fmt.Sprintf("%s/%s", folder, filename)

	// Upload file lên MinIO với metadata
	info, err := s.client.PutObject(
		context.Background(),
		s.config.MinIO.Bucket,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: s.getContentType(file.Filename), // Set MIME type
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Tạo presigned URL để access file (7 ngày)
	url, err := s.client.PresignedGetObject(
		context.Background(),
		s.config.MinIO.Bucket,
		objectName,
		time.Hour*24*7, // 7 days expiry
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &models.FileUploadResponse{
		Filename: objectName,
		URL:      url.String(),
		Size:     info.Size,
	}, nil
}

// UploadFileFromReader upload file từ io.Reader (dùng cho internal operations)
func (s *MinIOStorage) UploadFileFromReader(reader io.Reader, objectName string, size int64, contentType string) error {
	_, err := s.client.PutObject(
		context.Background(),
		s.config.MinIO.Bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	return err
}

// GetFile lấy file từ MinIO theo object name
// Trả về ReadCloser để stream file content
func (s *MinIOStorage) GetFile(objectName string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(
		context.Background(),
		s.config.MinIO.Bucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from MinIO: %w", err)
	}
	return object, nil
}

// DeleteFile xóa file từ MinIO
func (s *MinIOStorage) DeleteFile(objectName string) error {
	err := s.client.RemoveObject(
		context.Background(),
		s.config.MinIO.Bucket,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}
	return nil
}

// GetPresignedURL tạo presigned URL để access file
// URL có thời gian hết hạn để bảo mật
func (s *MinIOStorage) GetPresignedURL(objectName string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(
		context.Background(),
		s.config.MinIO.Bucket,
		objectName,
		expiry,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// ListFiles liệt kê tất cả files theo prefix
// Dùng để list files trong một folder
func (s *MinIOStorage) ListFiles(prefix string) ([]string, error) {
	ctx := context.Background()
	objectCh := s.client.ListObjects(ctx, s.config.MinIO.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true, // List recursively trong subfolders
	})

	var files []string
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		files = append(files, object.Key)
	}

	return files, nil
}

// generateUniqueFilename tạo filename unique bằng cách thêm timestamp
// Giúp tránh conflict khi upload files có tên giống nhau
func (s *MinIOStorage) generateUniqueFilename(originalFilename, folder string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// getContentType xác định MIME type dựa trên file extension
// Giúp browser hiểu cách xử lý file khi download
func (s *MinIOStorage) getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".epub":
		return "application/epub+zip"
	case ".mobi":
		return "application/x-mobipocket-ebook"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream" // Generic binary type
	}
}

// IsValidFileType kiểm tra file extension có được phép upload không
// Dựa trên whitelist trong config để bảo mật
func (s *MinIOStorage) IsValidFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	ext = strings.TrimPrefix(ext, ".") // Remove dot prefix
	
	// Kiểm tra extension có trong danh sách cho phép không
	for _, allowedType := range s.config.Upload.AllowedFileTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config là struct chính chứa tất cả cấu hình của ứng dụng
// Được chia thành các nhóm cấu hình riêng biệt để dễ quản lý
type Config struct {
	Server   ServerConfig  // Cấu hình server HTTP
	JWT      JWTConfig     // Cấu hình JSON Web Token
	MongoDB  MongoDBConfig // Cấu hình cơ sở dữ liệu MongoDB
	MinIO    MinIOConfig   // Cấu hình object storage MinIO
	Upload   UploadConfig  // Cấu hình upload file
}

// ServerConfig chứa cấu hình liên quan đến HTTP server
type ServerConfig struct {
	Port    string // Port để chạy server (default: 8080)
	GinMode string // Mode của Gin framework (debug/release)
}

// JWTConfig chứa cấu hình cho JSON Web Token
type JWTConfig struct {
	Secret      string // Secret key để sign/verify JWT token
	ExpiryHours int    // Thời gian hết hạn token (đơn vị: giờ)
}

// MongoDBConfig chứa cấu hình kết nối MongoDB
type MongoDBConfig struct {
	URI      string        // Connection string tới MongoDB
	Database string        // Tên database sử dụng
	Timeout  time.Duration // Timeout cho các operation
}

// MinIOConfig chứa cấu hình kết nối MinIO object storage
type MinIOConfig struct {
	Endpoint  string // Endpoint của MinIO server
	AccessKey string // Access key để authenticate
	SecretKey string // Secret key để authenticate
	UseSSL    bool   // Có sử dụng HTTPS không
	Bucket    string // Tên bucket mặc định để lưu file
}

// UploadConfig chứa cấu hình cho việc upload file
type UploadConfig struct {
	MaxFileSize      int64    // Kích thước file tối đa (bytes)
	AllowedFileTypes []string // Danh sách các loại file được phép upload
}

// Load đọc cấu hình từ file .env và environment variables
// Trả về pointer tới Config struct đã được khởi tạo với tất cả giá trị
func Load() *Config {
	// Cố gắng load file .env, nếu không có thì sử dụng environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Khởi tạo config với các giá trị default hoặc từ environment
	cfg := &Config{
		// Cấu hình server HTTP
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		// Cấu hình JWT authentication
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-change-me"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		// Cấu hình MongoDB database
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "volcanion_ebook"),
			Timeout:  time.Duration(getEnvAsInt("MONGODB_TIMEOUT", 10)) * time.Second,
		},
		// Cấu hình MinIO object storage
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
			Bucket:    getEnv("MINIO_BUCKET_NAME", "ebook-storage"),
		},
		// Cấu hình upload file
		Upload: UploadConfig{
			MaxFileSize:      parseFileSize(getEnv("MAX_FILE_SIZE", "100MB")),
			AllowedFileTypes: []string{"pdf", "epub", "mobi"},
		},
	}

	return cfg
}

// getEnv lấy giá trị từ environment variable với fallback value
// Trả về giá trị từ env nếu có, ngược lại trả về defaultValue
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt lấy giá trị integer từ environment variable
// Parse string thành int, nếu lỗi thì trả về defaultValue
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvAsBool lấy giá trị boolean từ environment variable
// Parse string thành bool, nếu lỗi thì trả về defaultValue
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// parseFileSize chuyển đổi string size thành bytes
// Hỗ trợ các đơn vị: KB, MB, GB (ví dụ: "100MB", "1GB")
func parseFileSize(size string) int64 {
	// Kiểm tra độ dài string tối thiểu
	if len(size) < 2 {
		return 100 * 1024 * 1024 // Default 100MB nếu format không hợp lệ
	}
	
	// Tách unit (2 ký tự cuối) và value (phần còn lại)
	unit := size[len(size)-2:]
	valueStr := size[:len(size)-2]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 100 * 1024 * 1024 // Default 100MB nếu parse lỗi
	}
	
	// Chuyển đổi theo đơn vị
	switch unit {
	case "KB":
		return value * 1024
	case "MB":
		return value * 1024 * 1024
	case "GB":
		return value * 1024 * 1024 * 1024
	default:
		return value // Trả về giá trị gốc nếu không có unit
	}
}

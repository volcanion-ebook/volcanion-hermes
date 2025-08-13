package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/volcanion/volcanion-hermes/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database struct chứa kết nối MongoDB và các thông tin cần thiết
type Database struct {
	Client   *mongo.Client   // MongoDB client để thực hiện operations
	DB       *mongo.Database // Database instance để truy cập collections
	Config   *config.Config  // App config chứa thông tin kết nối
}

// NewDatabase tạo kết nối mới tới MongoDB
// Thực hiện ping để test connection và tạo indexes cần thiết
func NewDatabase(cfg *config.Config) (*Database, error) {
	// Tạo context với timeout để tránh kết nối treo
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer cancel()

	// Tạo client options với connection string từ config
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test kết nối bằng cách ping tới MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Lấy database instance theo tên được config
	db := client.Database(cfg.MongoDB.Database)
	
	// Khởi tạo Database struct
	database := &Database{
		Client: client,
		DB:     db,
		Config: cfg,
	}

	// Tạo indexes để tối ưu hóa query performance
	// Nếu lỗi thì chỉ log warning, không fail toàn bộ application
	if err := database.createIndexes(ctx); err != nil {
		log.Printf("Warning: failed to create indexes: %v", err)
	}

	log.Println("Successfully connected to MongoDB")
	return database, nil
}

// createIndexes tạo các indexes cần thiết cho database
// Indexes giúp tăng tốc độ query và đảm bảo unique constraints
func (d *Database) createIndexes(ctx context.Context) error {
	// Tạo indexes cho users collection
	usersCollection := d.DB.Collection("users")
	userIndexes := []mongo.IndexModel{
		{
			// Index unique cho username để đảm bảo không trùng lặp
			Keys:    map[string]interface{}{"username": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			// Index unique cho email để đảm bảo không trùng lặp
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	
	// Thực hiện tạo indexes cho users collection
	if _, err := usersCollection.Indexes().CreateMany(ctx, userIndexes); err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	// Tạo indexes cho ebooks collection
	ebooksCollection := d.DB.Collection("ebooks")
	ebookIndexes := []mongo.IndexModel{
		{
			// Text index cho tìm kiếm full-text trên title, author, description
			Keys: map[string]interface{}{"title": "text", "author": "text", "description": "text"},
		},
		{
			// Index unique cho ISBN, sparse để cho phép null values
			Keys: map[string]interface{}{"isbn": 1},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			// Index theo thời gian tạo để sort theo mới nhất
			Keys: map[string]interface{}{"created_at": -1},
		},
		{
			// Index theo category để filter nhanh
			Keys: map[string]interface{}{"category": 1},
		},
		{
			// Index theo tags để tìm kiếm theo tag
			Keys: map[string]interface{}{"tags": 1},
		},
	}
	
	// Thực hiện tạo indexes cho ebooks collection
	if _, err := ebooksCollection.Indexes().CreateMany(ctx, ebookIndexes); err != nil {
		return fmt.Errorf("failed to create ebook indexes: %w", err)
	}

	return nil
}

// Close đóng kết nối MongoDB một cách graceful
// Sử dụng timeout để tránh treo application khi shutdown
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	return d.Client.Disconnect(ctx)
}

// GetCollection trả về MongoDB collection theo tên
// Wrapper method để dễ dàng truy cập collections
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.DB.Collection(name)
}

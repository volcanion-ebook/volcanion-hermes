package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/volcanion/volcanion-hermes/internal/database"
	"github.com/volcanion/volcanion-hermes/internal/models"
)

// UserService xử lý business logic liên quan đến user management
// Bao gồm authentication, user CRUD, và role management
type UserService struct {
	db         *database.Database // Database connection
	collection *mongo.Collection  // Users collection reference
}

// NewUserService tạo instance mới của UserService
func NewUserService(db *database.Database) *UserService {
	return &UserService{
		db:         db,
		collection: db.GetCollection("users"),
	}
}

// CreateUser tạo user mới trong hệ thống
// Hash password và set default roles nếu không được cung cấp
func (s *UserService) CreateUser(req *models.RegisterRequest) (*models.User, error) {
	// Hash password với bcrypt để bảo mật
	// Cost factor = DefaultCost (hiện tại là 10) cân bằng giữa security và performance
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default roles nếu không được cung cấp
	// Default role là "user" cho người dùng thông thường
	roles := req.Roles
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	// Tạo user object với thông tin đã được xử lý
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Roles:     roles,
		IsActive:  true, // Mặc định user được active
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user vào database
	result, err := s.collection.InsertOne(context.Background(), user)
	if err != nil {
		// Kiểm tra lỗi duplicate key (username hoặc email đã tồn tại)
		if mongo.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("username or email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Set ID từ result để trả về complete user object
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// GetUserByUsername tìm user theo username
// Sử dụng trong quá trình login authentication
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail tìm user theo email address
// Dùng cho password reset hoặc email verification
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByID tìm user theo ObjectID
// Sử dụng khi có ID từ JWT token hoặc relations
func (s *UserService) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// ValidatePassword so sánh plain password với hashed password
// Sử dụng bcrypt.CompareHashAndPassword để verify
func (s *UserService) ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// UpdateUser cập nhật thông tin user theo ID
// Tự động set updated_at timestamp
func (s *UserService) UpdateUser(id primitive.ObjectID, updates bson.M) error {
	// Thêm timestamp để track thời gian update
	updates["updated_at"] = time.Now()
	
	_, err := s.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	return err
}

// ListUsers lấy danh sách users với pagination
// Sort theo created_at descending (mới nhất trước)
func (s *UserService) ListUsers(page, limit int) ([]models.User, int64, error) {
	// Tính offset dựa trên page number
	skip := (page - 1) * limit
	
	// Đếm tổng số documents để tính pagination
	total, err := s.collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Tìm users với pagination và sorting
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
	cursor, err := s.collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(context.Background())

	// Decode tất cả results vào slice
	var users []models.User
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, 0, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, total, nil
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User đại diện cho một user trong hệ thống
// Chứa thông tin authentication và role-based access control
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username" validate:"required"`
	Email    string             `bson:"email" json:"email" validate:"required,email"`
	Password string             `bson:"password" json:"-"` // Không trả về password trong JSON response
	Roles    []string           `bson:"roles" json:"roles"` // Danh sách roles (admin, editor, user)
	IsActive bool               `bson:"is_active" json:"is_active"` // Trạng thái active/inactive
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

// Ebook đại diện cho một cuốn sách điện tử trong hệ thống
// Chứa metadata và thông tin file storage
type Ebook struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Author      string             `bson:"author" json:"author" validate:"required"`
	Publisher   string             `bson:"publisher" json:"publisher"`
	PublishYear int                `bson:"publish_year" json:"publish_year"`
	ISBN        string             `bson:"isbn" json:"isbn"` // International Standard Book Number
	Description string             `bson:"description" json:"description"`
	Language    string             `bson:"language" json:"language"`
	Category    string             `bson:"category" json:"category"` // Thể loại sách
	Tags        []string           `bson:"tags" json:"tags"`         // Tags để tìm kiếm và phân loại
	TotalPages  int                `bson:"total_pages" json:"total_pages"` // Tổng số trang
	FileSize    int64              `bson:"file_size" json:"file_size"`     // Kích thước file (bytes)
	FileFormat  string             `bson:"file_format" json:"file_format"` // Format file (PDF, EPUB, MOBI)
	CoverImage  string             `bson:"cover_image" json:"cover_image"` // Path tới ảnh bìa trong MinIO
	FilePath    string             `bson:"file_path" json:"file_path"`     // Path tới file ebook trong MinIO
	Pages       []EbookPage        `bson:"pages" json:"pages"`             // Danh sách các trang đã được chia nhỏ
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`   // ID của user tạo ebook
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// EbookPage đại diện cho một trang của ebook sau khi được chia nhỏ
// Mỗi trang có thể là một file PDF riêng biệt
type EbookPage struct {
	PageNumber int    `bson:"page_number" json:"page_number"` // Số trang
	FilePath   string `bson:"file_path" json:"file_path"`     // Path tới file trang trong MinIO
	FileSize   int64  `bson:"file_size" json:"file_size"`     // Kích thước file trang
	Text       string `bson:"text,omitempty" json:"text,omitempty"` // Text được OCR extract (optional)
}

// CreateEbookRequest là request body để tạo ebook mới
// Chứa các thông tin metadata cần thiết
type CreateEbookRequest struct {
	Title       string   `json:"title" validate:"required"`
	Author      string   `json:"author" validate:"required"`
	Publisher   string   `json:"publisher"`
	PublishYear int      `json:"publish_year"`
	ISBN        string   `json:"isbn"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// UpdateEbookRequest là request body để cập nhật ebook
// Sử dụng pointer để cho phép partial update (chỉ update field có giá trị)
type UpdateEbookRequest struct {
	Title       *string  `json:"title,omitempty"`
	Author      *string  `json:"author,omitempty"`
	Publisher   *string  `json:"publisher,omitempty"`
	PublishYear *int     `json:"publish_year,omitempty"`
	ISBN        *string  `json:"isbn,omitempty"`
	Description *string  `json:"description,omitempty"`
	Language    *string  `json:"language,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// LoginRequest là request body cho API đăng nhập
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest là request body cho API đăng ký tài khoản
type RegisterRequest struct {
	Username string   `json:"username" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"` // Tối thiểu 6 ký tự
	Roles    []string `json:"roles"` // Optional, default là ["user"]
}

// JWTClaims chứa thông tin claims trong JWT token
// Bao gồm thông tin user và RBAC data
type JWTClaims struct {
	UserID   string   `json:"user_id"`   // ID của user
	Username string   `json:"username"`  // Username
	Email    string   `json:"email"`     // Email
	Roles    []string `json:"roles"`     // Danh sách roles cho RBAC
	IsActive bool     `json:"is_active"` // Trạng thái active
	Exp      int64    `json:"exp"`       // Expiration timestamp
	Iat      int64    `json:"iat"`       // Issued at timestamp
}

// APIResponse là format chuẩn cho tất cả API responses
// Đảm bảo consistency trong việc trả về dữ liệu
type APIResponse struct {
	Success bool        `json:"success"`           // Trạng thái thành công/thất bại
	Message string      `json:"message"`           // Thông báo cho user
	Data    interface{} `json:"data,omitempty"`    // Dữ liệu trả về (nếu có)
	Error   string      `json:"error,omitempty"`   // Chi tiết lỗi (nếu có)
}

// PaginationRequest chứa thông tin phân trang từ query parameters
type PaginationRequest struct {
	Page  int `form:"page" json:"page"`   // Số trang (bắt đầu từ 1)
	Limit int `form:"limit" json:"limit"` // Số items per page
}

// EbookListResponse là response cho API list ebooks với phân trang
type EbookListResponse struct {
	Ebooks []Ebook `json:"ebooks"` // Danh sách ebooks trong trang hiện tại
	Total  int64   `json:"total"`  // Tổng số ebooks trong database
	Page   int     `json:"page"`   // Trang hiện tại
	Limit  int     `json:"limit"`  // Số items per page
}

// FileUploadResponse là response sau khi upload file thành công
type FileUploadResponse struct {
	Filename string `json:"filename"` // Tên file sau khi upload
	URL      string `json:"url"`      // URL để access file
	Size     int64  `json:"size"`     // Kích thước file (bytes)
}

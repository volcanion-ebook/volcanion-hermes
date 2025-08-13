package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/volcanion/volcanion-hermes/internal/config"
	"github.com/volcanion/volcanion-hermes/internal/models"
)

// JWTMiddleware xử lý authentication và authorization với JWT tokens
type JWTMiddleware struct {
	config *config.Config // App config chứa JWT secret và expiry
}

// NewJWTMiddleware tạo instance mới của JWTMiddleware
func NewJWTMiddleware(cfg *config.Config) *JWTMiddleware {
	return &JWTMiddleware{
		config: cfg,
	}
}

// GenerateToken tạo JWT token từ thông tin user
// Token chứa user info và roles cho RBAC
func (j *JWTMiddleware) GenerateToken(user *models.User) (string, error) {
	// Tạo claims với thông tin user và thời gian expired
	claims := &models.JWTClaims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles, // RBAC information
		IsActive: user.IsActive,
		Exp:      time.Now().Add(time.Hour * time.Duration(j.config.JWT.ExpiryHours)).Unix(),
		Iat:      time.Now().Unix(),
	}

	// Tạo token với signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
		"roles":    claims.Roles,
		"is_active": claims.IsActive,
		"exp":      claims.Exp,
		"iat":      claims.Iat,
	})

	// Sign token với secret key từ config
	tokenString, err := token.SignedString([]byte(j.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken xác thực và parse JWT token
// Trả về claims nếu token hợp lệ
func (j *JWTMiddleware) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	// Parse token với callback để verify signing method và secret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra signing method phải là HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// Trả về secret key để verify signature
		return []byte(j.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims từ token đã được verify
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Parse roles array từ interface{} sang []string
		roles := make([]string, 0)
		if rolesInterface, exists := claims["roles"]; exists {
			if rolesSlice, ok := rolesInterface.([]interface{}); ok {
				for _, role := range rolesSlice {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			}
		}

		// Tạo JWTClaims struct từ map claims
		return &models.JWTClaims{
			UserID:   claims["user_id"].(string),
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			Roles:    roles,
			IsActive: claims["is_active"].(bool),
			Exp:      int64(claims["exp"].(float64)),
			Iat:      int64(claims["iat"].(float64)),
		}, nil
	}

	return nil, jwt.ErrTokenInvalid
}

// AuthRequired là middleware yêu cầu authentication
// Kiểm tra JWT token trong Authorization header
func (j *JWTMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Authorization header required",
				Error:   "missing_auth_header",
			})
			c.Abort()
			return
		}

		// Extract token từ "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Nếu không có prefix "Bearer " thì format không đúng
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid authorization header format",
				Error:   "invalid_auth_format",
			})
			c.Abort()
			return
		}

		// Validate và parse JWT token
		claims, err := j.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid token",
				Error:   "invalid_token",
			})
			c.Abort()
			return
		}

		// Kiểm tra user có active không
		if !claims.IsActive {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "User account is inactive",
				Error:   "inactive_user",
			})
			c.Abort()
			return
		}

		// Kiểm tra token có hết hạn không
		if time.Now().Unix() > claims.Exp {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Token expired",
				Error:   "token_expired",
			})
			c.Abort()
			return
		}

		// Lưu user info vào context để các handler khác sử dụng
		c.Set("user", claims)
		c.Next()
	}
}

// RequireRoles là middleware kiểm tra RBAC permissions
// Yêu cầu user phải có ít nhất một trong các roles được chỉ định
func (j *JWTMiddleware) RequireRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy user info từ context (đã được set bởi AuthRequired middleware)
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "User not authenticated",
				Error:   "not_authenticated",
			})
			c.Abort()
			return
		}

		// Type assertion để lấy JWTClaims
		user, ok := userInterface.(*models.JWTClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Invalid user data",
				Error:   "invalid_user_data",
			})
			c.Abort()
			return
		}

		// Kiểm tra user có ít nhất một role required không
		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range user.Roles {
				if userRole == requiredRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		// Nếu không có quyền thì trả về Forbidden
		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Insufficient permissions",
				Error:   "insufficient_permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext helper để lấy user info từ Gin context
// Sử dụng trong các handlers để truy cập thông tin user hiện tại
func (j *JWTMiddleware) GetUserFromContext(c *gin.Context) (*models.JWTClaims, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, jwt.ErrTokenInvalid
	}

	user, ok := userInterface.(*models.JWTClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalid
	}

	return user, nil
}

# Volcanion Hermes - Ebook Management System

[![CI Status](https://github.com/your-org/volcanion-hermes/workflows/CI/badge.svg)](https://github.com/your-org/volcanion-hermes/actions)
[![Security Scan](https://github.com/your-org/volcanion-hermes/workflows/Security/badge.svg)](https://github.com/your-org/volcanion-hermes/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/volcanion-hermes)](https://goreportcard.com/report/github.com/your-org/volcanion-hermes)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Một hệ thống quản lý ebook hiện đại được xây dựng bằng **Golang** với kiến trúc microservices và các tính năng:

- 🔐 **Authentication & Authorization**: JWT với RBAC (Role-Based Access Control)
- 📚 **Ebook Management**: CRUD operations với MongoDB text search
- 🗃️ **File Storage**: MinIO S3-compatible object storage
- 🔍 **Search & Filter**: Full-text search và advanced filtering
- 📄 **File Upload**: Support PDF, EPUB, MOBI formats
- 🖼️ **Cover Images**: Image upload và presigned URLs
- 🚀 **Modern Stack**: Go 1.21, Gin, MongoDB, MinIO
- �️ **Security**: bcrypt password hashing, JWT expiration, input validation
- 📊 **Observability**: Health checks, structured logging
- 🐳 **Containerized**: Docker & Docker Compose support

## 🏗️ Kiến trúc hệ thống

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │    │   API Gateway   │    │   Load Balancer │
│  (Web, Mobile)  │◄───┤   (Nginx/...)   │◄───┤  (Optional)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                        ┌───────▼────────┐
                        │  Volcanion     │
                        │  Hermes API    │
                        │  (Gin Server)  │
                        └───────┬────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
        ┌───────▼────────┐ ┌───▼────┐ ┌────────▼────────┐
        │   MongoDB      │ │  MinIO │ │   Redis Cache   │
        │   (Database)   │ │ (Files)│ │   (Optional)    │
        └────────────────┘ └────────┘ └─────────────────┘
```

## 📁 Cấu trúc dự án

```
volcanion-hermes/
├── 📂 cmd/
│   └── server/
│       └── main.go              # 🚀 Entry point của ứng dụng
├── 📂 internal/
│   ├── config/
│   │   └── config.go           # ⚙️ Cấu hình ứng dụng & environment
│   ├── database/
│   │   └── database.go         # 🗄️ Kết nối MongoDB & indexes
│   ├── handlers/
│   │   ├── auth.go            # 🔐 HTTP handlers cho authentication
│   │   └── ebook.go           # 📚 HTTP handlers cho ebook CRUD
│   ├── middleware/
│   │   └── jwt.go             # 🛡️ JWT authentication & RBAC middleware
│   ├── models/
│   │   └── models.go          # 📋 Data models & API contracts
│   ├── services/
│   │   ├── user.go            # 👤 User business logic
│   │   └── ebook.go           # 📖 Ebook business logic
│   └── storage/
│       └── minio.go           # 🗃️ MinIO file storage service
├── 📂 .github/
│   ├── workflows/             # 🤖 GitHub Actions CI/CD
│   ├── dependabot.yml        # 🔄 Dependency auto-updates
│   └── README.md              # 📖 GitHub configuration docs
├── 📂 postman/
│   ├── Volcanion-Hermes-API.postman_collection.json  # 🧪 API testing
│   └── Volcanion-Hermes-Environment.postman_environment.json
├── 📂 docs/                   # 📚 Additional documentation
├── 🐳 Dockerfile              # Container definition
├── 🐳 docker-compose.yml      # Multi-service setup
├── ⚙️ .env.example            # Environment template
├── 📦 go.mod                  # Go module dependencies
├── 🔧 Makefile               # Build automation
└── 📝 README.md              # This file
```

## 🚀 Quick Start

### ⚡ Với Docker (Khuyến nghị)

```bash
# 1. Clone repository
git clone https://github.com/your-org/volcanion-hermes.git
cd volcanion-hermes

# 2. Copy environment file
cp .env.example .env

# 3. Start tất cả services với Docker Compose
docker-compose up -d

# 4. Kiểm tra health
curl http://localhost:8080/health

# 5. Import Postman collection để test API
# File: postman/Volcanion-Hermes-API.postman_collection.json
```

### 🔧 Development Setup

```bash
# 1. Cài đặt Go dependencies
go mod download

# 2. Start MongoDB và MinIO (hoặc dùng Docker)
docker-compose up -d mongodb minio

# 3. Setup environment
cp .env.example .env
# Chỉnh sửa .env file theo cấu hình local

# 4. Run application
make run
# hoặc: go run cmd/server/main.go
```

## 📋 Yêu cầu hệ thống

| Component | Version | Purpose |
|-----------|---------|---------|
| **Go** | 1.21+ | Backend language |
| **MongoDB** | 7.0+ | Database |
| **MinIO** | Latest | Object storage |
| **Docker** | 20.10+ | Containerization |
| **Git** | 2.30+ | Version control |

### Optional Tools
- **Postman** - API testing
- **MongoDB Compass** - Database GUI
- **MinIO Console** - Storage management

## 🛠️ Installation & Setup

### 📥 Method 1: Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/your-org/volcanion-hermes.git
cd volcanion-hermes

# Setup environment
cp .env.example .env

# Tùy chỉnh .env file nếu cần
nano .env

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f volcanion-hermes

# Stop services
docker-compose down
```

### 🔧 Method 2: Manual Setup

#### 1. Install Dependencies

```bash
# Install Go (Ubuntu/Debian)
wget -c https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install MongoDB
curl -fsSL https://pgp.mongodb.com/server-7.0.asc | sudo gpg --dearmor -o /usr/share/keyrings/mongodb-server-7.0.gpg
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt-get update
sudo apt-get install -y mongodb-org

# Install MinIO
wget https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio
sudo mv minio /usr/local/bin/
```

#### 2. Setup Services

```bash
# Start MongoDB
sudo systemctl start mongod
sudo systemctl enable mongod

# Start MinIO
mkdir -p ~/minio-data
minio server ~/minio-data --console-address ":9001" &

# Verify services
mongosh --eval "db.adminCommand('ping')"
curl http://localhost:9000/minio/health/live
```

#### 3. Application Setup

```bash
# Clone và setup
git clone https://github.com/your-org/volcanion-hermes.git
cd volcanion-hermes

# Install Go dependencies
go mod download

# Setup environment
cp .env.example .env
# Chỉnh sửa .env theo cấu hình local

# Build application
make build

# Run application
make run
```

## ⚙️ Configuration

### Environment Variables

```bash
# Server Configuration
PORT=8080                    # HTTP server port
GIN_MODE=release             # Gin mode: debug, release, test

# JWT Configuration
JWT_SECRET=your-secret-key   # JWT signing secret (thay đổi trong production)
JWT_EXPIRATION=24h          # Token expiration time

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017/volcanion_hermes
MONGODB_DATABASE=volcanion_hermes
MONGODB_TIMEOUT=30s

# MinIO Configuration
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET=volcanion-hermes

# Application Settings
MAX_FILE_SIZE=50MB          # Max upload file size
ALLOWED_FILE_TYPES=pdf,epub,mobi
LOG_LEVEL=info              # Log level: debug, info, warn, error
```

### 🏗️ Build Options

```bash
# Development build
make build

# Production build với optimizations
make build-prod

# Cross-platform builds
make build-all

# Docker build
make docker-build

# Clean builds
make clean
```

## 📚 API Documentation

### 🔐 Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Đăng ký user mới | ❌ |
| POST | `/api/v1/auth/login` | Đăng nhập | ❌ |
| GET | `/api/v1/auth/profile` | Lấy profile user | ✅ |
| POST | `/api/v1/auth/refresh` | Refresh JWT token | ✅ |

### 📖 Ebook Endpoints

| Method | Endpoint | Description | Auth Required | Roles |
|--------|----------|-------------|---------------|-------|
| GET | `/api/v1/ebooks` | Danh sách ebooks | ❌ | All |
| GET | `/api/v1/ebooks/search` | Tìm kiếm ebooks | ❌ | All |
| GET | `/api/v1/ebooks/:id` | Chi tiết ebook | ❌ | All |
| GET | `/api/v1/ebooks/:id/download` | Download ebook | ✅ | All |
| POST | `/api/v1/ebooks` | Tạo ebook mới | ✅ | Admin, Editor |
| PUT | `/api/v1/ebooks/:id` | Cập nhật ebook | ✅ | Admin, Editor |
| DELETE | `/api/v1/ebooks/:id` | Xóa ebook | ✅ | Admin, Editor |
| POST | `/api/v1/ebooks/:id/upload` | Upload file ebook | ✅ | Admin, Editor |
| POST | `/api/v1/ebooks/:id/cover` | Upload cover image | ✅ | Admin, Editor |

### 🏥 System Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |

## 🧪 Testing với Postman

### Import Collection

1. **Download Postman**: https://www.postman.com/downloads/
2. **Import Collection**:
   - File: `postman/Volcanion-Hermes-API.postman_collection.json`
   - Click "Import" trong Postman
   - Chọn file collection

3. **Import Environment**:
   - File: `postman/Volcanion-Hermes-Environment.postman_environment.json`
   - Click "Import" → "Environment"
   - Chọn file environment

### Test Workflow

```bash
# 1. Health Check
GET /health

# 2. Register User
POST /api/v1/auth/register
{
  "username": "testuser",
  "email": "test@example.com", 
  "password": "SecurePassword123!",
  "full_name": "Test User",
  "role": "user"
}

# 3. Login (auto-saves token)
POST /api/v1/auth/login
{
  "username": "testuser",
  "password": "SecurePassword123!"
}

# 4. Create Ebook (needs admin/editor role)
POST /api/v1/ebooks
{
  "title": "Go Programming Guide",
  "author": "John Doe",
  "description": "Comprehensive Go guide"
}

# 5. Upload File
POST /api/v1/ebooks/{id}/upload
Form-data: file=ebook.pdf

# 6. Search & List
GET /api/v1/ebooks?page=1&limit=10
GET /api/v1/ebooks/search?q=golang
```

### Auto-Generated Tests

Collection bao gồm các test tự động:
- ✅ Status code validation
- ✅ Response structure validation  
- ✅ Token management
- ✅ Performance checks
- ✅ Error handling

## 🔨 Development

### 🏃‍♂️ Running Tests

```bash
# Unit tests
make test

# Tests với coverage
make test-coverage

# Integration tests (yêu cầu MongoDB + MinIO)
make test-integration

# Benchmark tests
make benchmark

# Race condition tests
make test-race
```

### 🔍 Code Quality

```bash
# Linting
make lint

# Format code
make fmt

# Security scan
make security

# Vulnerability check
make vuln-check

# All quality checks
make quality
```

### 📊 Monitoring & Debugging

```bash
# Health check
curl http://localhost:8080/health

# Logs
docker-compose logs -f volcanion-hermes

# Performance profiling
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profiling  
go tool pprof http://localhost:8080/debug/pprof/heap
```

## 🚀 Deployment

### 🐳 Docker Deployment

```bash
# Build production image
docker build -t volcanion-hermes:latest .

# Run with production environment
docker run -d \
  --name volcanion-hermes \
  -p 8080:8080 \
  --env-file .env.production \
  volcanion-hermes:latest
```

### ☁️ Cloud Deployment

#### AWS ECS/EKS
```bash
# Build và push image
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin <account>.dkr.ecr.us-west-2.amazonaws.com
docker tag volcanion-hermes:latest <account>.dkr.ecr.us-west-2.amazonaws.com/volcanion-hermes:latest
docker push <account>.dkr.ecr.us-west-2.amazonaws.com/volcanion-hermes:latest
```

#### Google Cloud Run
```bash
# Build và deploy
gcloud builds submit --tag gcr.io/<project>/volcanion-hermes
gcloud run deploy volcanion-hermes --image gcr.io/<project>/volcanion-hermes --platform managed
```

#### Digital Ocean App Platform
```yaml
# .do/app.yaml
name: volcanion-hermes
services:
- name: api
  source_dir: /
  github:
    repo: your-org/volcanion-hermes
    branch: main
  dockerfile_path: Dockerfile
  http_port: 8080
  instance_count: 1
  instance_size_slug: basic-xxs
  env:
  - key: PORT
    value: "8080"
```

### 🔄 CI/CD Pipeline

Project sử dụng GitHub Actions cho CI/CD:

- **CI**: Lint, test, security scan, build
- **CD**: Auto-deploy tới staging/production
- **Security**: Dependency scanning, secret detection
- **Release**: Automated releases với semantic versioning

Chi tiết xem: [.github/README.md](.github/README.md)

## 🏗️ Architecture & Design

### 🎯 Design Principles

- **Clean Architecture**: Separation of concerns
- **SOLID Principles**: Maintainable code
- **Domain-Driven Design**: Business logic focus
- **API-First**: Contract-driven development
- **Security by Design**: Defense in depth

### 🔐 Security Features

- **JWT Authentication**: Stateless authentication
- **RBAC Authorization**: Role-based access control
- **Password Hashing**: bcrypt với salt
- **Input Validation**: Comprehensive validation
- **File Type Validation**: Safe file uploads
- **Rate Limiting**: API protection (optional)
- **CORS**: Cross-origin protection
- **Security Headers**: HTTP security headers

### 📈 Performance Features

- **Connection Pooling**: Database connections
- **Concurrent Processing**: Goroutines
- **Efficient Queries**: MongoDB indexes
- **File Streaming**: Large file handling
- **Presigned URLs**: Direct S3 access
- **Caching**: Redis support (optional)

## 🤝 Contributing

### 🔧 Development Workflow

1. **Fork** repository
2. **Create** feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** changes: `git commit -m 'Add amazing feature'`
4. **Push** branch: `git push origin feature/amazing-feature`
5. **Open** Pull Request

### 📋 Code Standards

- **Go Conventions**: Follow effective Go guidelines
- **Comments**: Document public APIs
- **Tests**: Minimum 80% coverage
- **Commit Messages**: Conventional commits format
- **PR Template**: Use provided template

### 🐛 Issue Reporting

Template cho bug reports:
```markdown
**Bug Description**
Clear description của issue

**Steps to Reproduce**
1. Step one
2. Step two
3. See error

**Expected Behavior**
What should happen

**Environment**
- OS: [Ubuntu 22.04]
- Go Version: [1.21.0]
- Docker: [20.10.0]
```

## 📄 License

Project này được license dưới [MIT License](LICENSE).

```
MIT License

Copyright (c) 2025 Volcanion Hermes

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 🙏 Acknowledgments

- **[Gin Web Framework](https://gin-gonic.com/)** - HTTP web framework
- **[MongoDB Go Driver](https://go.mongodb.org/mongo-driver)** - MongoDB integration  
- **[MinIO Go SDK](https://min.io/docs/minio/linux/developers/go/minio-go.html)** - Object storage
- **[JWT-Go](https://github.com/golang-jwt/jwt)** - JWT implementation
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Password hashing

## 📞 Support & Contact

- **Documentation**: [Wiki](https://github.com/your-org/volcanion-hermes/wiki)
- **Issues**: [GitHub Issues](https://github.com/your-org/volcanion-hermes/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/volcanion-hermes/discussions)
- **Email**: support@your-org.com

---

<div align="center">

**⭐ Star project nếu bạn thấy hữu ích!**

[![GitHub stars](https://img.shields.io/github/stars/your-org/volcanion-hermes?style=social)](https://github.com/your-org/volcanion-hermes/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/your-org/volcanion-hermes?style=social)](https://github.com/your-org/volcanion-hermes/network/members)

**Made with ❤️ by [Your Team Name]**

</div>
  -p 9000:9000 \
  -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"

# Hoặc tải về binary
# https://docs.min.io/docs/minio-quickstart-guide.html
```

### 6. Khởi động ứng dụng

```bash
go run cmd/server/main.go
```

Ứng dụng sẽ chạy tại `http://localhost:8080`

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Đăng ký tài khoản | No |
| POST | `/api/v1/auth/login` | Đăng nhập | No |
| GET | `/api/v1/auth/profile` | Lấy thông tin profile | Yes |
| POST | `/api/v1/auth/refresh` | Refresh token | Yes |

### Ebooks

| Method | Endpoint | Description | Auth Required | Roles |
|--------|----------|-------------|---------------|--------|
| GET | `/api/v1/ebooks` | Danh sách ebook | No | - |
| GET | `/api/v1/ebooks/search` | Tìm kiếm ebook | No | - |
| GET | `/api/v1/ebooks/:id` | Chi tiết ebook | No | - |
| GET | `/api/v1/ebooks/:id/download` | Tải ebook | Yes | - |
| POST | `/api/v1/ebooks` | Tạo ebook mới | Yes | admin, editor |
| PUT | `/api/v1/ebooks/:id` | Cập nhật ebook | Yes | admin, editor |
| DELETE | `/api/v1/ebooks/:id` | Xóa ebook | Yes | admin, editor |
| POST | `/api/v1/ebooks/:id/upload` | Upload file ebook | Yes | admin, editor |
| POST | `/api/v1/ebooks/:id/cover` | Upload ảnh bìa | Yes | admin, editor |

## Roles và Permissions

- **admin**: Toàn quyền trên hệ thống
- **editor**: Quản lý ebook (tạo, sửa, xóa)
- **user**: Xem và tải ebook

## Ví dụ sử dụng API

### 1. Đăng ký tài khoản

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123",
    "roles": ["admin"]
  }'
```

### 2. Đăng nhập

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

### 3. Tạo ebook mới

```bash
curl -X POST http://localhost:8080/api/v1/ebooks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "Golang Programming",
    "author": "John Doe",
    "publisher": "Tech Books",
    "publish_year": 2023,
    "description": "A comprehensive guide to Go programming",
    "category": "Technology",
    "language": "English",
    "tags": ["golang", "programming", "backend"]
  }'
```

### 4. Upload file ebook

```bash
curl -X POST http://localhost:8080/api/v1/ebooks/<ebook-id>/upload \
  -H "Authorization: Bearer <your-token>" \
  -F "file=@/path/to/your/ebook.pdf"
```

### 5. Tìm kiếm ebook

```bash
curl "http://localhost:8080/api/v1/ebooks/search?q=golang&page=1&limit=10"
```

## Cấu trúc Database

### Users Collection

```json
{
  "_id": "ObjectId",
  "username": "string",
  "email": "string",
  "password": "string (hashed)",
  "roles": ["string"],
  "is_active": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Ebooks Collection

```json
{
  "_id": "ObjectId",
  "title": "string",
  "author": "string",
  "publisher": "string",
  "publish_year": "number",
  "isbn": "string",
  "description": "string",
  "language": "string",
  "category": "string",
  "tags": ["string"],
  "total_pages": "number",
  "file_size": "number",
  "file_format": "string",
  "cover_image": "string",
  "file_path": "string",
  "pages": [
    {
      "page_number": "number",
      "file_path": "string",
      "file_size": "number",
      "text": "string"
    }
  ],
  "created_by": "ObjectId",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Tính năng

### 🔐 Authentication & Authorization
- JWT tokens với RBAC payload
- Role-based access control
- Token refresh mechanism
- Password hashing với bcrypt

### 📚 Ebook Management
- CRUD operations cho ebook
- Upload file và ảnh bìa
- Metadata management
- Full-text search
- Categorization và tagging

### 🗃️ File Storage
- MinIO integration
- Presigned URLs cho download
- File type validation
- Automatic file organization

### 🔍 Search & Filter
- Text search trên title, author, description
- Filter theo category, author
- Pagination support

## Development

### Build ứng dụng

```bash
go build -o bin/server cmd/server/main.go
```

### Chạy tests

```bash
go test ./...
```

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/server .
COPY --from=builder /app/.env .
CMD ["./server"]
```

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

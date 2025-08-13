# API Examples for Volcanion Hermes

This directory contains example API requests for testing the Volcanion Hermes ebook management system.

## Prerequisites

- The server should be running on `http://localhost:8080`
- MongoDB and MinIO should be running
- You can use tools like curl, Postman, or HTTPie

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication Flow

### 1. Register a new user

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

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

Response will include a JWT token. Save it for subsequent requests.

### 3. Get user profile

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Ebook Management

### 1. Create a new ebook

```bash
curl -X POST http://localhost:8080/api/v1/ebooks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Learning Golang",
    "author": "John Doe",
    "publisher": "Tech Publications",
    "publish_year": 2023,
    "isbn": "978-1234567890",
    "description": "A comprehensive guide to learning Go programming language",
    "language": "English",
    "category": "Programming",
    "tags": ["golang", "programming", "backend", "tutorial"]
  }'
```

### 2. Get all ebooks

```bash
curl -X GET "http://localhost:8080/api/v1/ebooks?page=1&limit=10"
```

### 3. Get ebook by ID

```bash
curl -X GET http://localhost:8080/api/v1/ebooks/EBOOK_ID
```

### 4. Update ebook

```bash
curl -X PUT http://localhost:8080/api/v1/ebooks/EBOOK_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Advanced Golang Programming",
    "description": "An updated comprehensive guide to advanced Go programming"
  }'
```

### 5. Search ebooks

```bash
curl -X GET "http://localhost:8080/api/v1/ebooks/search?q=golang&page=1&limit=5"
```

### 6. Filter ebooks by category

```bash
curl -X GET "http://localhost:8080/api/v1/ebooks?category=Programming&page=1&limit=10"
```

### 7. Filter ebooks by author

```bash
curl -X GET "http://localhost:8080/api/v1/ebooks?author=John&page=1&limit=10"
```

## File Upload

### 1. Upload ebook file

```bash
curl -X POST http://localhost:8080/api/v1/ebooks/EBOOK_ID/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/book.pdf"
```

### 2. Upload cover image

```bash
curl -X POST http://localhost:8080/api/v1/ebooks/EBOOK_ID/cover \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "cover=@/path/to/cover.jpg"
```

### 3. Download ebook

```bash
curl -X GET http://localhost:8080/api/v1/ebooks/EBOOK_ID/download \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

This returns a JSON response with a presigned download URL.

## Error Handling

All API responses follow this format:

```json
{
  "success": true/false,
  "message": "Description of the result",
  "data": {},
  "error": "Error details (if any)"
}
```

Common HTTP status codes:
- `200`: Success
- `201`: Created
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `500`: Internal Server Error

## Environment Variables

Make sure your `.env` file is properly configured:

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=volcanion_ebook
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
JWT_SECRET=your-secret-key
```
